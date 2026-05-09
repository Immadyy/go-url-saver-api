⭐ 1. In-Memory Storage vs Persistent Database

What I did was use an in-memory slice with a mutex lock.

The reason I did this was because it was the way to get a working prototype.

I did not have to deal with zero dependencies.

There was no Docker, no connection strings and no migration headaches while I focused on structure.


Tradeoffs are:

* setup, no external dependencies

* Blazing reads and writes

* All data is gone on server restart

* No instance support each server has its own data

* No query capabilities, no filtering, pagination or search


What I am changing is moving to PostgreSQL.

Java is different from Go.

In Java Spring Boot with JPA/Hibernate gives you interfaces that auto-generate queries.

In Go you have to write SQL by hand unless you pull in an ORM like GORM.

The tradeoff is that Go gives you control and Spring gives you more speed once you know it.


⭐ 2. Interface-Based Storage

What I did was define a LinkStore interface. Made Database implement it.

LinkService depends on the interface, not the type.

The reason I did this was because it is dependency injection.

If I swap the in-memory DB for PostgreSQL I change zero lines in LinkService.

I just write a struct that satisfies LinkStore and pass it in.


Tradeoffs are:

* Testable. I can mock LinkStore in unit tests

* Swappable backends without touching business logic

* More boilerplate than calling a global DB variable

* Go interfaces are implicit so I need discipline to keep them clean


Java is different from Go.

In Java Spring does this automatically with annotations.

In Go you wire everything manually in main.

Gos way is more explicit and debuggable.

Springs way is faster to write but hides the wiring.

In an interview saying I prefer wiring for smaller projects but Springs DI container scales better for large teams shows mature thinking.


⭐ 3. Graceful Shutdown

What I did was signal.NotifyContext catches SIGTERM/SIGINT, server.Shutdown with a 10-second timeout drains requests before closing.

The reason I did this was because without this killing the server drops in-flight requests.

Users get errors.

With this ongoing requests finish then the server stops cleanly.


Tradeoffs are:

* Production-ready behavior

* Respects Kubernetes/Docker stop signals

* Adds complexity to a main function

* The 10-second timeout is arbitrary. If a request takes longer it still gets killed


Java is different from Go.

In Java Spring Boot has this built-in via spring.lifecycle.timeout-per-shutdown-phase.

You configure it in a properties file.

In Go I have to write it myself. I understand exactly whats happening.


4. Custom Middleware vs Third-Party Router

What I did was use Go 1.22s http.NewServeMux with method-based routing. Wrote my logging middleware by wrapping http.ResponseWriter.

The reason I did this was because Go 1.22 finally supports method-based routing in the library.

Before this everyone used gorilla/mux or chi.

I wanted to try the native approach.

The custom middleware captures HTTP status codes for logging.


Tradeoffs are:

* Zero dependencies for routing

* Full control over middleware behavior

* No built-in path parameters. I used query params instead

* No middleware chaining library. Have to nest if I add more middleware


Java is different from Go.

In Java Spring Boot has annotations and middleware is done via Filters or Interceptors.

Springs approach is annotation-driven and declarative.

Gos approach is procedural and explicit.

Neither is better it's about team preference.


⭐ 5. URL Validation

What I did was in ValidateLink if a URL doesn't have http:// or https:// I prepend https:// automatically.

Then I validate with url.ParseRequestURI.

The reason I did this was because users are lazy.

They'll type google.com instead of https://google.com.

so instead of rejecting it I fix it for them.

This is a UX decision, not technical one.


Tradeoffs are:

* Better user experience

* Reduces support/bug reports about URL saving

* Assumes HTTPS. May break for ftp:// or other protocols

* Silent correction might confuse power users who made a typo


Java is the same as Go in this case.

The interesting comparison is: in Go this validation lives in the service layer as a function.

In Spring you could do it as a validator or even a custom Jackson deserializer.

Different placement, concept.


6. API Response Wrapper

What I did was every response uses APIResponse structure.

The reason I did this was because it is a response format.

Frontend always knows to look for message and data.

Errors from the service layer get wrapped in HTTP error codes instead.


Tradeoffs are:

* client-side parsing

* Easy to add metadata later

* Inconsistent. GetAllHandler returns raw array not wrapped in APIResponse

* Some APIs prefer success, error, data format


Fix: Wrap GetAllHandler response

I noticed this while writing this doc.

Thanks, past me.


⭐ 7. Mutex Locking on In-Memory Store

What I did was every database method calls d.mu.Lock and d.mu.Unlock.

The reason I did this was because Go servers handle each request in its goroutine.

Without the mutex two simultaneous writes could corrupt the slice.

The go run -race flag would catch this.


Tradeoffs are:

* Safe for use

* Simple to implement

* Lock contention. If 1000 requests hit at once they all queue up

* GetAll locks too. Reads block writes


Java is different from Go.

In Java Spring Boot with JPA handles concurrency via database transactions and isolation levels.

Gos explicit locking forces you to think about concurrency earlier.

In Java many devs ignore it until production breaks.

This is a point to make in an interview. In Go concurrency is explicit and I had to think about it from day one.


8. Error Handling Style

What I did was every function returns result and error.

No exceptions.

Errors are values.

I check them immediately with if err != nil.


Tradeoffs are:

* Explicit. You know which functions can fail

* No hidden control flow like try/catch

* Verbose. Lots of if err != nil blocks

* Easy to accidentally ignore an error by not checking


Java is different from Go.

In Java Java uses checked exceptions and unchecked exceptions.

Spring Boot wraps many checked exceptions into runtime ones.

Gos approach is simpler and more explicit. Javas exception hierarchy gives you more granularity.

Interview gold: I like Gos forced error handling because Junior devs can't hide from errors. Javas exception types make debugging production incidents faster.


⭐ 9. Single-File vs Package Structure

What I did was everything in main.go initially.

The reason I did this was because it was iteration.

No import cycles.

Easy to share with others for feedback.


Tradeoffs are:

* Fast to prototype

* To paste/share for code review

* 300+ lines in one file gets unwieldy

* No clear boundaries between layers


What I am changing is splitting into files.

Java is different from Go.

In Java Java enforces package structure from the start.

Spring Boot projects have a layout.

Go gives you freedom. Which's great for small projects and dangerous for large ones.

Another good interview point.