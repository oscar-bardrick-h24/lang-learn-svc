# Please Read

I'm very sorry but I didnt quite manage to finish this in time. Most of the key functionality is there but some things are very rough around the edges. I started a new job last week so as you can probably imagine I've unfortunately been pretty overloaded and didnt manage to find enough time to properly do the task justice. I am very keen on working at Babel though so I hope that you'll still consider me despite the shortcomings of this piece of work. I wrote this very quickly and there are, as you will see, some issues which I simply didn't have time to address. I just hope I will get a chance to discuss those in a follow-up interview.

 In the past I've primarily worked with and am most comfortable with MySQL, so I originally started building the service around a MySQL backend. However midway through I switched to Postgres since the task description mentioned it was preferred and I wanted to exploit the ARRAY type. After converting the code to be Postgres compatible I realised that the `github.com/lib/pq` driver does not support advanced Postgres types like ARRAY and that I would have to use `github.com/jackc/pgx`. The issue with that was that pgx has a much different interface which would've required bigger changes that I wasn't sure whether I'd have time to complete. As such the final design is a bit of a mish mash between the original MySQL design and a Postgres design using the JSON type for the courses lesson_list. Hopefully you will give me the opportunity to discuss what is wrong with this and also some of the added flexibility it potentially delivers.

 Although there wasnt time to fully cover the module in tests, I've tried to write a reasonable amount to show that:
  1. the app is written in a manner that makes testing easy
  2. I know how to test and i understand the importance of testing

  One of he largest shortcomings is that the error handling is not comprehensive or consistent. I realise this, unfortunately I didnt get the time to fix this.

  ## Testing

  unit tests can be run with `make unit-test`
  integration tests can be run with `make integration-test` once `docker-compose up` has succeeded and containers are ready for action

  ## Design

  I always try follow the community standards for package structuring/naming as set out in [this extremely useful resource](https://github.com/golang-standards/project-layout)

  I have always admired [Matt Ryer's style](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html) of writing Go web apps for it's simplicity and cleanliness. It's a design that I believe takes the principle of KISS seriously and makes writing, testing and maintaining web-services a pleasure. This is the design I use to guide me.

  the General idea is:
  - pkg main: main.go that collects app configuration from the enviroment and uses those values to construct app dependencies and inject them into the App type
  - pkg app: App type encapsulates all dependencies required to run the application. It's sole responsibility is HTTP request/response. This means that the app package is made up of middlewares, handlers and routing. Each handler aims to do nothing beyond receiving data from requests, handing that data off to a Service and finally constructing an appopriate response based on Service's return value. we dont want the handler doing anything beyond request/response logic, as this would be a violation of the single-responsibilty principle
  - pkg domain: the domain package is home to all business logic types and methods. This is where the interesting stuff happens. Once the handlers have extracted the pertinent bits of data from the request, domain.Service is the component that actually understands and processes that data
  - pkg repo: below the Service layer is the data layer that handles everything persistence related.

  Each layer defines the underlying layer as an interface so testing is nice and easy.


Many thanks for reading and I hope to hear back soon, please don't hesitate to reach out with any questions. my email is oscarjo@tuta.io