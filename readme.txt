Timeout details:
    - client timeout
    - server timeout
    - timeout at the load balancers.

Go provides following ways to create request with timeouts:
    - http.client
        - Which encompasses the whole request cycle from 'Dial' to 'Response Body'.

    - context
        - 'context' package provides useful tools to handle timeout, 
            Deadline and cancellable Requests via 'WithTimeout', 'WithDeadline' and 'WithCancel' methods. 
            Using 'WithTimeout', you can add the timeout to the http.Request using req.WithContext method.

    - http.Transport
        -You can specify the timeout also using the low level implementation of creating a custom 'http.Transport' 
            using 'DialContext' and use it to create the 'http.client'

SetDeadline:
    The network primitive that Go exposes to implement timeouts: Deadlines.
    Deadlines are an absolute time which when reached makes all I/O operations fail with a timeout error.

    keep in mind that all timeouts are implemented in terms of Deadlines, so they do NOT reset every time data is sent or received.
        - net.Conn().SetDeadline(time.Now().Add(time.Second * 30))

    // Server
    srv := &http.Server{
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    log.Println(srv.ListenAndServe())

    You should set both timeouts when you deal with untrusted clients and/or networks, 
    so that a client can't hold up a connection by being slow to write or read.

client Timeout Detail:
    For more granular control, there are a number of other more specific timeouts you can set:
    - 'net.Dialer.Timeout' limits the time spent establishing a TCP connection (if a new one is needed).
    - 'http.Transport.TLSHandshakeTimeout' limits the time spent performing the TLS handshake.
    - 'http.Transport.ResponseHeaderTimeout' limits the time spent reading the headers of the response.
    - 'http.Transport.ExpectContinueTimeout' limits the time the client will wait between sending the request headers when including an Expect: 100-continue and receiving the go-ahead to send the body.

    Ex:
        c := &http.Client{
            Transport: &http.Transport{
                Dial: (&net.Dialer{
                        Timeout:   30 * time.Second,
                        KeepAlive: 30 * time.Second,
                }).Dial,
                TLSHandshakeTimeout:   10 * time.Second,
                ResponseHeaderTimeout: 10 * time.Second,
                ExpectContinueTimeout: 1 * time.Second,
            }
        }



Ref:
https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
https://itnext.io/http-request-timeouts-in-go-for-beginners-fe6445137c90