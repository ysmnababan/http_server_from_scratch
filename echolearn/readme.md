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

- [ ] What a radix tree is and why echo uses it instead of a plain map
- [ ] How route parameters like `:id` are stored and extracted from the tree
- [ ] How wildcard routes `*` differ from param routes `:param` in the tree
- [ ] What happens when two routes conflict — how echo resolves priority
- [ ] How `e.Group()` works and what it actually stores vs a standalone route
- [ ] How echo handles `405 Method Not Allowed` vs `404 Not Found` — what triggers each

### Practice questions

**Q4.** You register `/users/:id` and `/users/me`. A request comes in for `/users/me`.
Which route matches and why?

> Your answer:

---

**Q5.** What is the structural difference between how echo stores `e.GET("/a/b")` vs
`e.GET("/a/:param")` in the radix tree?

> Your answer:

---

**Q6.** If you register only `POST /users` and a client sends `GET /users`, does echo
return 404 or 405? Trace exactly which code produces that response.

> Your answer:

---

## Phase 3 — Middleware chain

**Files to read:** `middleware/logger.go`, `middleware/recover.go`, `echo.go` (`Use`, `Group`)

### Concepts

- [ ] How `MiddlewareFunc` wraps `HandlerFunc` — the exact type signature and what it means
- [ ] How `e.Use()` middleware differs from group-level and route-level middleware in execution order
- [ ] What happens if a middleware does not call `next(c)`
- [ ] How echo's `Recover` middleware catches panics — what it does with the stack trace
- [ ] How `Logger` middleware captures response status code — why it needs to wrap the `ResponseWriter`

### Practice questions

**Q7.** You have global middleware A, group middleware B, and route middleware C. What is
the exact execution order? Does the order of `e.Use()` calls matter?

> Your answer:

---

**Q8.** Why can't the `Logger` middleware just read `c.Response().Status` after calling
`next(c)`? What problem does this create and how does echo solve it?

> Your answer:

---

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
