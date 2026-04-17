# Echo Framework — Learning Plan & Contribution Readiness Checklist

A structured plan to go from echo user to echo contributor. Complete each phase in order.
Check off concepts as you understand them by reading the source. Answer every practice question
before moving to the next phase.

**Source to read:** <https://github.com/labstack/echo>

**Progress:** \_\_\_/ 24 concepts

---

## Phase 1 — Echo core: entry point & context

**Files to read:** `echo.go`, `context.go`

### Concepts

- [x] How `echo.New()` initializes the struct and what fields it sets up
      it initilize handler, logger, tls, router, middleware, and even `http.Server`
- [x] Why `Echo` implements `http.Handler` and what `ServeHTTP` does first
      so it can mimic the same handler as the http handler and run the servehttp.
      `serverhttp()` is then mainly used to fetch the preregistered `handler` for
      that specific resource path. if it is not, then the handler will be executed
      by the `defaultServeMux` and not the echo's handler.
      the serveHTTP itself take the `echo.context` first from the `sync.Pool`
- [x] What the `sync.Pool` is doing in echo and why it's used for `Context`
      `echo.context` is a struct that is use for each handler later on. so it encasulates
      the `*http.request` and wrapped `http.responsewriter`. also it has the `echo.echo`
      itself. this context is quite big, so to make it more performant, the `echo.context`
      can be fetched using `sync.pool` to reduce the reallocation.
- [x] The lifecycle of a `Context` — acquired, used, reset, returned to pool
      the context is acquired from the pool at the beginnin of the `ServeHTTP` method,
      then it will be reset and populated with the request and resnpose, then find the
      router. and the router will be executed to get the handler. then apply middleware,
      and the registered handler that uses that `context`. after that, the context will be
      return to the pool.
- [x] What `c.Param()`, `c.QueryParam()`, and `c.Bind()` each do under the hood
      `c.queryparam` fetch the query param using the underlying `url.url` inside the
      http request.
      `c.param` uses their own param store and value which i still dont totaly understand
      now.
      c.bind is general method that bind the request body to a struct. it will check the
      path param first, then query param, then parse the body based on the mime type (it
      can be json, xml, multipart, etc).
      actualy the bind uses the `defaultBinder` which is initialized when calling the
      `echo.New()` method.
- [x] How `c.JSON()` writes the response — what headers it sets and in what order
      `c.json` writes the response by settting the content type to `application/json`
      and then set the status code.
      inside `context`, there is a `JSONSerializer` which is a echo's own interface for
      json serialization. by default, it wraps the default `json.Marshal`.
      using this serializer, the `c.JSON` serialize the data and write it to response body
      through the `context.Response()` which is a wrapped `http.ResponseWriter`.
- [x] What `DefaultHTTPErrorHandler` does and when it's called
      echo has a field called `HTTPErrorHandler` which is a function for handling the error.
      it is called when the registered handler return and error. this error will be handler
      using something like this : `HTTPErrorHandler(err,c)`. inside it, it will check if
      the error is echo's own sentinel error or not. if not, then it will converted to
      500 internal server error. then the response will be written using `c.json()`.
      So the `DefaultHTTPErrorHandler` is just a default implementation of this error handler,
      which is created when calling `echo.New()`. you can replace it with your own
- [x] How returning an `error` from a handler flows through to the error handler
      like explained above, there are globale handler for handling error which is default to
      `DefaultHTTPErrorHandler`. firstly, the handler will check if the response is already
      committed or not (via `c.json` or other method). if not, then it will check the error
      type if it is echo's own error or not. if not, then return internal server error.
      then call the `c.json` again to finalize

question:

- why echo struct use context (not directly but through sync.pool) and context struct
  also use echo struct? it is like circular dependency
-

### Practice questions

**Q1.** If you call `c.JSON(200, data)` and then `c.JSON(400, err)` in the same handler,
what happens? Why?

> Your answer:
> it will produce 2 response because the json encoder will just write the the message to the
> response body without checking it. but it is different if for the status code. it will check
> if it is committed or not. if it is commited, then it will just ignore the new status code.
> basically calling `c.JON`will commit the response, so the second call will just write to the
> body without changing the status code. so the final result is 200 with the body of
> both `err` and `data`.

---

**Q2.** Why does echo use `sync.Pool` for `Context` instead of just creating a new one
per request?

> Your answer:
> it is because the echo context is quite big, so to make it more perfromat, it can reuse the
> existing context from the pool. and to make it safe to use, it uses the `sync` package and
> reset the context before using it

---

**Q3.** What is the difference between returning `echo.NewHTTPError(404, "not found")`
vs returning a plain `errors.New("not found")` from a handler?

> Your answer:
> assume that we use default error handler, the first error will return as not found
> because it is a `echo.HTTPErro`. so it will return as is.
> but for the latter because it is just plain error and not `echo.HTTPError`, it will send as
> internal server error so the error format will be the same

---

## Phase 2 — Router internals

**Files to read:** `router.go`

### Concepts

- [x] What a radix tree is and why echo uses it instead of a plain map
      _Ans_: radix tree is a trie-based data structure that instead of storing a single char or
      'data' inside a node, it purposefuly store the common prefix from a set of data to
      reduce memory and increase time of access. instead of searching all the node -- which can
      costly -- the complexity can be reduced to O(m) where `m` is the length of a string.
      echo uses it because the way the request path is constructed. it consist of dynamic key
      like `:id` that can't be easily stored in key-value data. also, radix tree can improve
      both performance and memory especially with server that has a lot of endpoints.
- [x] How route parameters like `:id` are stored and extracted from the tree
      _Ans_: `:id` are stored as `pnames`, that will be added when adding the handler along with
      the handler itself. both stored in struct called `routemethod`. and the actual value for
      the `:id` itself will be populated when the endpoint is being hit. the server will eventually
      run the `ServeHTTP` method which call the method `router#Find`. this method will find the
      handler based on the path(and method). at the same time, this method will populate the
      `context` from pool and add the `pnames` and `pvalues` pair. if the handler later on want to
      get the param, it can quickly search from `pnames` and `pvalues` inside the `context`.
- [x] How wildcard routes `*` differ from param routes `:param` in the tree
      _Ans_: Basically wildcard routes is route where the path after the `*` symbol will be redirected
      to the registered handler on the `*` node. if the wildcard is found then it will go to
      the handler even if there is remaining path. the behaviour is slightly different with
      the `:param` because it can still have the next fragment after the `:param` symbol.
      When registered, param can create up to 3 new nodes, 1 for path before `:`, 1 for path
      until the `:` symbol, and the rest is for the path after the `:`.
      but different with the wildcard routes, where it only create 2 nodes with without the
      preceding path.
- [x] What happens when two routes conflict — how echo resolves priority
      _Ans_: when registering the same route twice, it will rewrite the handler into the latest
      handler. when trying to get the handler (when request is received), it will find the
      exact match first.
- [x] How `e.Group()` works and what it actually stores vs a standalone route
      _Ans_: `e.Group()` basically creates a new struct called `group` which is a wrapper around
      the `e` itself with added prefix (and middleware). the returned `group` has its own register
      http method. inside it, it will call `echo.add` but the path will be enriched with the
      group prefix.
- [x] How echo handles `405 Method Not Allowed` vs `404 Not Found` — what triggers each
      _Ans_: method not allowed is sent whenever there are matching url but the method is not
      match. but the `not found` is sent when there is no matching url.

### Practice questions

**Q4.** You register `/users/:id` and `/users/me`. A request comes in for `/users/me`.
Which route matches and why?

> Your answer:
> these paths will eventually introduce 3 nodes with something like this :

```
/user/  => :
        => me
```

so when the request `/users/me` comes, it will use the `m` as a label to find the
`me` node. so it will find the `/users/me` first because it has identical match.
the node kind priority is `staticKind` > `paramKind` > `catchAllKind`.
so it will find the static node first before

---

**Q5.** What is the structural difference between how echo stores `e.GET("/a/b")` vs
`e.GET("/a/:param")` in the radix tree?

> Your answer:
> `staticKind` will likely combined in one single node if possible because echo uses
> LCP as their algorithm. so as long as no other path that has common prefix, the
> the path will be stored in single node. but this is different with the `:param`.
> because it will automatically created at least 2 nodes where the lcp is started
> before the `:` symbol. so `:` will be the new label for the root to use whenever
> it wants to check its children

---

**Q6.** If you register only `POST /users` and a client sends `GET /users`, does echo
return 404 or 405? Trace exactly which code produces that response.

> Your answer:
> The router will return 405 because there is matching path but not with the method

---

## Phase 3 — Middleware chain

**Files to read:** `middleware/logger.go`, `middleware/recover.go`, `echo.go` (`Use`, `Group`)

### Concepts

- [x] How `MiddlewareFunc` wraps `HandlerFunc` — the exact type signature and what it means
      _Ans_: `MiddlewareFunc` is function who receives a `HandlerFunc` and return a `HandlerFunc`.
      it means that the function `next` can be executed after some process and that's is the
      point of using middleware. techically speaking, it can also handle after the `next`
      execution. so it basically means that the middleware will act like a onion wrapper
      that wrap a particular handler with before and after execution.
      it look like this

  ```
    func Middleware(next HandlerFunc) HandlerFunc {
      return func(c echo.Context) error {
        // do something here
        return next(c)
      }
    }

    mw1
     mw2
      handler
     mw2
    mw1
  ```

- [x] How `e.Use()` middleware differs from group-level and route-level middleware in execution order
      _Ans_: `e.Use()` middleware is a global middleware, so it will applied to all the handler.
      but the group level only applied to endpoint which belong to the group and the
      route-level middlware will only applied to a specific endpoint
- [x] What happens if a middleware does not call `next(c)`
      _Ans_:It will not execute the handler inside it. meaning that the output for and endpoint will likely
      be empty, or as if the handler is called but the function is empty.
- [x] How echo's `Recover` middleware catches panics — what it does with the stack trace
      _Ans_: recover middleware has a defer function inside it where it catch all the panic from the
      `next(c)` execution result. it will check the return type of the `panic` and also show the error
      stack. The `Recover` middleware also has the `logerrorfunc` which can be added using the config.
      The stack can be displayed or not depend on the configuration too. But it has its own default
      configuration.
- [x] How `Logger` middleware captures response status code — why it needs to wrap the `ResponseWriter`
      The logger capture response status code by accessing the c.response(). There is no wrapping the
      response writer, but the `c.Response()` is a wrapper around the `http.ResponseWriter` that also
      has a field for status code. so when the handler write to the response, it will update the
      status code in the `c.Response()`. so the logger can access this field after calling
      `next(c)` to get the final status code.

### Practice questions

**Q7.** You have global middleware A, group middleware B, and route middleware C. What is
the exact execution order? Does the order of `e.Use()` calls matter?

> Your answer:
> the execution order is A>B>C. The order of e.Use calls matter because it will apply
> the middleware by wrapping the handler from the last to the first. so if the request come, it will
> execute A first because it is the most outer wrapper.

---

**Q8.** Why can't the `Logger` middleware just read `c.Response().Status` after calling
`next(c)`? What problem does this create and how does echo solve it?

> Your answer:
> it can read the status code after calling `next(c)`, but the problem is that the status code is not
> set until the handler write to the response. so if the handler just return without writing to the response,
> then the status code will be 0. echo solve this by wrapping the `http.ResponseWriter` with its own `Response`
> struct that has a field for status code. so whenever the handler write to the response, it will update the
> status code in the `Response` struct. so the logger can access this field after calling `next(c)` to
> get the final status code

## Phase 4 — Binder, renderer & extensibility

**Files to read:** `bind.go`, `echo.go` (`Validator`, `Renderer`)

### Concepts

- [ ] How `c.Bind()` decides whether to decode JSON, form, or query params — what it checks and in what order
- [ ] How echo uses Go reflection in its default binder to populate struct fields
- [ ] How to replace the default binder with a custom one and when you'd want to
- [ ] What `Validator` is in echo and why echo doesn't ship a default one
- [ ] How `c.Render()` works and what interface a custom renderer must implement

### Practice questions

**Q9.** A request comes in with `Content-Type: application/json` but also has query params.
What does `c.Bind()` populate and what does it ignore?

> Your answer:

---

**Q10.** Why does echo not ship with a built-in validator despite shipping with a
built-in binder?

> Your answer:

---

## Contribution readiness threshold

You are ready to contribute when you can answer all 10 questions without looking at the
source, and can locate the relevant code for each answer within 2 minutes of opening
the repo.

At that point, go back to the issue tracker and read open issues again. You will read
them completely differently.

### Suggested first contribution targets

Once you pass the threshold, look for issues in these areas — they match what you now
understand deeply:

- Router edge cases — overlapping params, wildcard conflicts
- Binder bugs — unexpected behavior with specific `Content-Type` combinations
- Middleware correctness — execution order edge cases with groups
- Missing test coverage — behaviors that exist but aren't tested

---

## Resources

- Echo source: <https://github.com/labstack/echo>
- Echo docs: <https://echo.labstack.com>
- Go `net/http` source (for reference): <https://github.com/golang/go/tree/master/src/net/http>
- Your HTTP from scratch notes: `README.md`
