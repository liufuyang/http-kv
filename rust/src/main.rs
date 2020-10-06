use hyper::service::{make_service_fn, service_fn};
use hyper::{Body, Method, Request, Response, Server, StatusCode};

use lazy_static::lazy_static;
use dashmap::DashMap;
use std::sync::Arc;
use std::thread;
use std::time::{Duration, SystemTime};

use prometheus::{exponential_buckets, register_gauge, register_histogram_vec, Encoder, TextEncoder, HistogramVec, Gauge};

const EXPIRE_DURATION: Duration = Duration::from_secs(60);

lazy_static! {
    // Cache solution
    // https://docs.rs/dashmap/3.11.10/dashmap/struct.DashMap.html
    static ref DASHMAP: Arc<DashMap<String, Value>> = {
        let m =  DashMap::new();
        let arc = Arc::new(m);
        let map = arc.clone();

        thread::spawn(move || {
            let vacuum_cycle_sleep = if EXPIRE_DURATION < Duration::from_secs(10) {EXPIRE_DURATION} else {Duration::from_secs(10)};
            loop {
                println!("Removing keys");
                map.retain(|_, v: &mut Value| !v.expired());
                METRIC_CACHE_SIZE.set(map.len() as f64);
                thread::sleep(vacuum_cycle_sleep);
            }
        });

        arc
    };

    // Metrics
    // https://gist.github.com/breeswish/bb10bccd13a7fe332ef534ff0306ceb5
    static ref METRIC_REQUEST_DURATION: HistogramVec = register_histogram_vec!(
        "request_duration_seconds",
        "Histogram of HTTP request duration in seconds",
        &["method"],
        exponential_buckets(0.005, 2.0, 10).unwrap()
        ).unwrap();

    static ref METRIC_CACHE_SIZE: Gauge = register_gauge!(
        "cache_size",
        "Size of the cache"
        ).unwrap();
}

struct Value {
    s: String,
    time: SystemTime,
}

impl Value {
    pub fn new(s: String) -> Value {
        Value { s, time: SystemTime::now() }
    }

    pub fn to_string(&self) -> String {
        self.s.clone()
    }

    pub fn expired(&self) -> bool {
        SystemTime::now().duration_since(self.time)
            .map_or(true, |d| d > EXPIRE_DURATION)
    }
}

async fn kv_handler(req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
    let start_time = SystemTime::now();
    let result = match req.method() {
        // Serve some instructions at /
        &Method::GET => {
            let key: &str = req.uri().path();
            if key.len() < 2 {
                return Ok(Response::builder().status(400).body(Body::from("Must provide a key in the path")).unwrap());
            }
            if key.to_lowercase().eq("/size") {
                return Ok(Response::new(Body::from(DASHMAP.len().to_string())));
            }
            if key.to_lowercase().eq("/metrics") {
                let encoder = TextEncoder::new();
                let mut buffer = vec![];
                let mf = prometheus::gather();
                encoder.encode(&mf, &mut buffer).unwrap();
                return Ok(Response::builder()
                    .header(hyper::header::CONTENT_TYPE, encoder.format_type())
                    .body(Body::from(buffer))
                    .unwrap());
            }
            let result = match DASHMAP.get(key) {
                Some(v) => Ok(Response::new(Body::from((*v).to_string()))),
                None => Ok(Response::default())
            };
            (result, "get")
        }

        &Method::POST => {
            let key: String = req.uri().path().to_string();
            if key.len() < 2 || key.to_lowercase().eq("/size") {
                return Ok(Response::builder().status(400).body(Body::from("Must provide a key in the path")).unwrap());
            }
            let bytes = hyper::body::to_bytes(req.into_body()).await?;
            DASHMAP.insert(key, Value::new(String::from_utf8(bytes.to_vec()).unwrap()));
            METRIC_CACHE_SIZE.inc();

            (Ok(Response::default()), "post")
        }

        // Return the 404 Not Found for other routes.
        _ => {
            let mut not_found = Response::default();
            *not_found.status_mut() = StatusCode::NOT_FOUND;
            (Ok(not_found), "other")
        }
    };

    let duration = SystemTime::now().duration_since(start_time).unwrap();
    METRIC_REQUEST_DURATION
        .with_label_values(&[result.1])
        .observe(duration.as_secs_f64());
    return result.0;
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let addr = ([127, 0, 0, 1], 8082).into();

    let service = make_service_fn(|_| async { Ok::<_, hyper::Error>(service_fn(kv_handler)) });

    let server = Server::bind(&addr).serve(service);

    println!("Listening on http://{}", addr);

    server.await?;

    Ok(())
}