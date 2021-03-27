# go-runs

Things to do:

- [x] Receive an error message that contains, name, message, tags
- [x] Store a set of resolutions to handle an error
- [x] Each resolution contains a set of steps (a markdown document) 
- [x] Each error triggers a session where operators can check what they executed and mark down results
- [ ] Sessions should be saved to DB
- [ ] Create multiple projects 
- [ ] Login and RBAC
- [ ] Git as a Markdown backend

TODOs
* <del>proper http status codes</del>
* <del>remove usage of struct pointers in favor of interfaces</del>
* <del>use immutables (hide fields) where it makes sense</del>
* Add more tests
* Fix TODOs
* Add proper go docs
* Create docker-compose to run