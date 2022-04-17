---
theme: gaia
_class: lead
paginate: true
backgroundColor: #fff
backgroundImage: url('https://marp.app/assets/hero-background.svg')
marp: true
---

# **Assignment**

Simple Rest API implementation

---

# Technologies Used

- PostgreSQL as the database system
- Golang/ Java ( SpringBoot, Apache Camel)
- Docker, Docker compose, Kubernetes for containerization and hosting infrastructure
- Minikube for setting up/ testing Kubernetes configurations
- Cmake to automate build tasks in Go source

---

# Design Decisions

## SpringBoot/ Apache Camel App
- Connection Pooling for better performance
- Use streaming and parallel executions in bulk message processing

--- 

## Go Application
- Layered architecture to make testable code
- Hold dependancies in a custom DI container (I have not used any third party DI containers since this is a very small app)
- Not using ORM for 3 reasons (There are some good ORM libraries present for GO)
    - Need control on low level stuff so we can write better optimizations
    - Can write and maintain SQL stuff in code since this is a small application
    - ORM comes with an added complexity and performance cost
- JSON logging to faciliate faster log parsing and indexing
- auto tagging builds/ images based on git tags
---

## TODO

01. Improve test coverage (Configure pipelines to run tests automatically and report coverage %)
02. Expose application matrices in prometheus format
03. improvements in logging
04. Write benchmarks on the hot code path
05. Automate docker builds on release tags
06. Store sensitive data in Secrets (databse passwords/usernames)
07. propagate correlation IDs through boundaries and write it in logs
08. Add/configure service health probes