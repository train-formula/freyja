vcl 4.0;

backend default {
	.host = "0.0.0.0";
	.port = "${FREYJA_PORT}";
}

backend s3USWest2 {
	.host = "s3-us-west-2.amazonaws.com";
	.port = "80";
}

backend linkerd {
	.host = "localhost";
	.port = "14080";
}

acl freyja {
	"localhost";
  	"127.0.0.1";
  	"::1";
}

sub vcl_hit {

}

sub vcl_recv {

	if (req.method != "GET" &&
      req.method != "HEAD" &&
      req.method != "PUT" &&
      req.method != "POST" &&
      req.method != "TRACE" &&
      req.method != "OPTIONS" &&
      req.method != "PATCH" &&
      req.method != "DELETE") {
	    /* Non-RFC2616 or CONNECT which is weird. */
	    /*Why send the packet upstream, while the visitor is using a non-valid HTTP method? */
	    return(synth(404, "Non-valid HTTP method!"));
  	}

	if ( req.method == "POST" || req.method == "PUT" || req.method == "DELETE" || req.method == "CONNECT" || req.method == "OPTIONS" || req.method == "TRACE" || req.method == "PATCH") {
		return (synth(403,""));
	}


	if( req.http.host ~ "s3-us-west-2.amazonaws.com" ){

		if (client.ip ~ freyja) {

			if(req.http.Use-Linkerd) {
				set req.backend_hint = linkerd;
			} else {
				set req.backend_hint = s3USWest2;
			}

		} else {
			return (synth(403,""));
		}
	}
}