package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict Transport Security", "max-age=63072000; includeSubDomains;preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Powered-By", "Django") // Puts people on the wrong path regarding tech used.

		//Explore these
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
	})
}

// basic middleware format
/*

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){})
}

*/
/*

dns prefetching is a technique used to resolve domain names before they are needed, reducing latency when the browser needs to connect to those domains. Turning it off reduces risk of dns based attacks.
X-Frame_Options is a security header that prevents clickjacking attacks by controlling whether a page can be displayed in a frame or iframe.
X-XSS-Protection is a security header that enables the browser's built-in cross-site scripting (XSS) filter to prevent XSS attacks.
X-Content-Type-Options is a security header that prevents browsers from interpreting files as a different MIME type than what is specified in the Content-Type header, mitigating MIME sniffing attacks.
Strict-Transport-Security is a security header that enforces secure (HTTPS) connections to the server. Prevents man in the middle attacks by ensuring that browsers only connect to the server over HTTPS.
Content-Security-Policy is a security header that helps prevent various types of attacks, such as cross-site scripting (XSS) and data injection attacks, by specifying which content sources are allowed to be loaded by the browser.
Referrer-Policy is a security header that controls how much referrer information is included with requests made from the page. Setting it to "no-referrer" prevents any referrer information from being sent, enhancing privacy and security.

*/

/*
space for additional security headers
Server header is used to identify the software used by the server. Setting it to an empty string can help obscure the technology stack, making it harder for attackers to target specific vulnerabilities.
X-Permitted-Cross-Domain-Policies is a security header that controls how cross-domain policies are handled, preventing unauthorized access to resources.
Cache-Control is a security header that controls caching behavior, preventing sensitive information from being cached by browsers or intermediaries.
Cross-Origin-Resource-Policy is a security header that controls how resources can be shared across different origins, enhancing security by preventing unauthorized access to resources.
Cross-Origin-Opener-Policy is a security header that controls how a document can interact with other documents from different origins, preventing cross-origin attacks.
Cross-Origin-Embedder-Policy is a security header that controls how a document can embed resources from different origins, enhancing security by preventing unauthorized access to embedded resources.
Access-Control-Allow-Headers is a security header that specifies which headers can be used in cross-origin requests, enhancing security by preventing unauthorized headers from being sent.
Access-Control-Allow-Methods is a security header that specifies which HTTP methods are allowed in cross-origin requests, enhancing security by preventing unauthorized methods from being used.
Access-Control-Allow-Credentials is a security header that indicates whether credentials (such as cookies or HTTP authentication) can be included in cross-origin requests, enhancing security by controlling credential sharing.
Permissions-Policy is a security header that allows or denies the use of certain browser features, such as geolocation or microphone access,
enhancing security by controlling which features can be used by the page.

*/
