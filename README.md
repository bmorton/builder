# Builder

A conventional Docker image builder that simply accepts webhooks from any Github repository, builds an image for that repository, and pushes it to the supplied registry tagged with its git SHA using the same name as the repo.

## Examples

* Triggering a build

  ```ShellSession
  $ curl http://localhost:3000/webhooks/github -H "Content-Type: application/json" -d @support/github_sample.json -v
  * Hostname was NOT found in DNS cache
  *   Trying ::1...
  * Connected to localhost (::1) port 3000 (#0)
  > POST /webhooks/github HTTP/1.1
  > User-Agent: curl/7.37.1
  > Host: localhost:3000
  > Accept: */*
  > Content-Type: application/json
  > Content-Length: 7408
  > Expect: 100-continue
  >
  < HTTP/1.1 100 Continue
  < HTTP/1.1 200 OK
  < Content-Type: application/json; charset=utf-8
  < Date: Wed, 01 Apr 2015 10:05:10 GMT
  < Content-Length: 233
  <
  {"id":"319db017-21ab-4ac0-be26-4bfd8688b7a5","repository_name":"hello-world","clone_url":"https://github.com/deployster/hello-world.git","commit_id":"96ac589c5a5d4366446e6675598bed1671913521","git_ref":"refs/heads/refactor-polling"}
  * Connection #0 to host localhost left intact
  ```

* Listing running builds

  ```ShellSession
  $ curl http://localhost:3000/builds -v
  * Hostname was NOT found in DNS cache
  *   Trying ::1...
  * Connected to localhost (::1) port 3000 (#0)
  > GET /builds HTTP/1.1
  > User-Agent: curl/7.37.1
  > Host: localhost:3000
  > Accept: */*
  >
  < HTTP/1.1 200 OK
  < Content-Type: application/json; charset=utf-8
  < Date: Wed, 01 Apr 2015 10:05:16 GMT
  < Content-Length: 235
  <
  [{"id":"319db017-21ab-4ac0-be26-4bfd8688b7a5","repository_name":"hello-world","clone_url":"https://github.com/deployster/hello-world.git","commit_id":"96ac589c5a5d4366446e6675598bed1671913521","git_ref":"refs/heads/refactor-polling"}]
  * Connection #0 to host localhost left intact
  ```

* Following the progress of a build

  ```ShellSession
  $ curl http://localhost:3000/builds/319db017-21ab-4ac0-be26-4bfd8688b7a5 -v
  * Hostname was NOT found in DNS cache
  *   Trying ::1...
  * Connected to localhost (::1) port 3000 (#0)
  > GET /builds/319db017-21ab-4ac0-be26-4bfd8688b7a5 HTTP/1.1
  > User-Agent: curl/7.37.1
  > Host: localhost:3000
  > Accept: */*
  >
  < HTTP/1.1 200 OK
  < Content-Type: text/plain
  < Date: Wed, 01 Apr 2015 10:05:24 GMT
  < Transfer-Encoding: chunked
  <
  Step 0 : FROM yammer/ruby:2.2.0
   ---> 7828aa49f4ac
  Step 1 : MAINTAINER Brian Morton "brian@mmm.hm"
   ---> Using cache
   ---> dc0f4b99910c
  Step 2 : ADD Gemfile /app/Gemfile
   ---> Using cache
   ---> 7fc9532741d6
  Step 3 : ADD Gemfile.lock /app/Gemfile.lock
   ---> Using cache
   ---> bf0446b9c8c6
  Step 4 : RUN bash -l -c "cd /app && bundle"
   ---> Using cache
   ---> 44ad8af43010
  ...
  Pushing [=========================>                         ]    512 B/1.024 kB 0
  Pushing [==================================================>] 1.024 kB/1.024 kB
  Pushing [==================================================>] 1.024 kB/1.024 kB
  Image successfully pushed
  Pushing tag for rev [626431e69f0d] on {http://192.168.59.103:5000/v1/repositories/hello-world/tags/96ac589}
  ```
