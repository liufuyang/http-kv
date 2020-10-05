use hyper::service::{make_service_fn, service_fn};
use hyper::{Body, Method, Request, Response, Server, StatusCode};

use lazy_static::lazy_static;
use dashmap::DashMap;
use std::sync::Arc;

lazy_static! {
    static ref DASHMAP: Arc<DashMap<String, String>> = {
        let m =  DashMap::new();
        m.insert("k1".to_string(), "foo".to_string());
        m.insert("k2".to_string(), "bar".to_string());
        m.insert("k3".to_string(), "baz".to_string());
        Arc::new(m)
    };
}

async fn kv_handler(req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
    match req.method() {
        // Serve some instructions at /
        &Method::GET => {
            let paths: Vec<&str> = req.uri().path().splitn(3, '/').filter(|v| !v.is_empty()).collect();
            if paths.len() < 1 {
                return Ok(Response::builder().status(400).body(Body::from("Must provide a key in the path")).unwrap());
            }

            let key = paths.get(0).unwrap();
            match DASHMAP.get(*key) {
                Some(v) => Ok(Response::new(Body::from((&*v).clone()))),
                None => Ok(Response::default())
            }
        }

        &Method::POST => {
            let paths: Vec<&str> = req.uri().path().splitn(3, '/').filter(|v| !v.is_empty()).collect();
            if paths.len() < 1 {
                return Ok(Response::builder().status(400).body(Body::from("Must provide a key in the path")).unwrap());
            }

            let key = paths.get(0).unwrap().to_string();
            let bytes = hyper::body::to_bytes(req.into_body()).await?;
            DASHMAP.insert(key, String::from_utf8(bytes.to_vec()).unwrap());

            Ok(Response::default())
        }

        // Return the 404 Not Found for other routes.
        _ => {
            let mut not_found = Response::default();
            *not_found.status_mut() = StatusCode::NOT_FOUND;
            Ok(not_found)
        }
    }
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