![Demo](http://g.recordit.co/XwYmeWdePJ.gif)
# CoderSchool Golang Course - URL Shortener

1. **Submitted by:** Arthur Nguyen
2. **Time spent:** 2 days

## Set of User Stories

### Requirements
#### The list of redirection should be maintained in a command line tool, what can:

* [x] Manipulate YAML config file. Where the redirection list peristently stored.
* [x] Implement append to the list: urlshorten configure -a dogs -u www.dogs.com
* [x] Implement remove from the list: urlshorten -d dogs
* [x] List redirections: urlshorten -l
* [x] Run HTTP server on a given port: urlshorten run -p 8080
* [x] Prints usage info: urlshorten -h

### As a bonus exercises you can also...

* [x] Track number of times each redirection is used. When the user uses urlshorten -l, the user should see redirections ranked by how many times they have been used.
* [x] Provide a default shortening, if no example is given. For example, if dogs is not given, generate a random alphanumeric string of length 8.
* [ ] Build a Handler that doesn't read from a map but instead reads from a database. Whether you use BoltDB, SQL, or something else is entirely up to you.
