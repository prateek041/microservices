---
title: "Domain Boundaries"
description: ""
date: 2025-05-26
---

## Why Kubernetes

If we are building our micro-services with the notion of not deploying it on
Kubernetes, we would have to go extra mile to implement the following things,

- Service Discovery.
- Load balancing.
- Orchestration and Management.
- Configuration management and Environment Variables.
- Health checks.

And many more. Our business logic and implementation remains the same
irrespective of the deployment platform only the **Infrastructure Glue**
will change.

## Handling Scaling

When defining our micro-services, we should consider how each service will need
to scale independently, For example:

- **Product Catalog**: Might experience high read traffic but relatively fewer write
  operations. We should design it to scale horizontally (more replicas) efficiently
  for reads.
- **Order Management**: Could have bursts of write activity during peak ordering
  times likes sales and might benefit from different scaling characteristics
  than the product catalog. [TODO: What kind of Scaling?].
- **User Management**: Might have high read traffic during login but fewer user
  creation.
- **Payment Processing**: Resource intensive but short lived.
- **Notification Service**: Might need to handle a large volume of asynchronous messages.

Since we are building this application with the notion of building it for Kubernetes
type environment, it's networking patterns might differ.

## Communication patterns in a Kubernetes environment

- **Internal Communication via Kubernetes Services**: Primary way for services to
  communicate with each other is through **services**. Each service gets a stable
  DNS name (e.g., `product-catalog-service`) and a virtual IP address.

  Kubernetes handles the routing of traffic to the underlying pods, providing
  basic load balancing. This simplifies service discovery and makes internal
  communication more reliable. We will primarily use synchronous HTTP or gRPC
  calls via these Service names for request/response interactions.

- **External Access via Ingress or LoadBalancer**: To expose our application
  to the outside world (e.g., user browsers or mobile apps), we'll typically use
  an Ingress controller or a LoadBalancer service. These provide external access
  points and can handle routing of external requests to the appropriate internal
  services (often via an API Gateway).

- **Asynchronous Communication via Message Queues**: For decoupled communication
  and background processing, we'll deploy a message broker like RabbitMQ within
  our Kubernetes cluster. Our services can then send and receive messages via
  Kubernetes Services that expose the message broker.

## Designing for Failure and Independent deployability

In a distributed system, failures are inevitable. We need to design our
`microservices` to be resilient and handle failures gracefully:

- **Statelessness**: Designing services to be stateless (not storing any
  persistent data within the service instance) makes them easier to scale and
  replace in case of failure. Any persistent data should be stored in external
  backing services (databases, caches).

- **Health Checks**: Implementing health check endpoints (/health) in each service
  allows Kubernetes to monitor their status and restart failing Pods. Readiness
  probes ensure that a Pod is ready to serve traffic before it's added to the
  Service's endpoints. They are very important because if implemented wrong, they
  can cut off the entire service from the rest of the application.

- **Timeouts and Retries**: When one service calls another, we should configure
  appropriate timeouts to prevent indefinite blocking. For transient failures,
  implementing retry mechanisms with backoff can improve resilience.

- **Circuit Breakers**: To prevent cascading failures, we can use circuit
  breaker patterns. If a service is failing consistently, the calling service
  can stop making requests for a period, giving the failing service time to recover.

## Communication Styles (Kubernetes Native)

Since we are working in a Kubernetes environment, let's properly define how our
microservices will communicate with each other in more details.

### Synchronous Communication via Kubernetes Services

Our Go micro-services can make standard HTTP or gRPC calls using these DNS names
as the `hostname` in the URL. For example, to call the Product Catalog service from
the User Management service, we would construct the URL like
`http://product-catalog-service:8080/products/some-id` (assuming the Product Catalog
service is named product-catalog-service and listens on port 8080 within its Pods).
Kubernetes will resolve `product-catalog-service` to the appropriate IP address and
load balance across the available Pods.

### Asynchronour communication using message queues (e.g., RabbitMQ deployed on K8s)

- **Deploying a Message Broker on Kubernetes**: For asynchronous communication, we'll
  deploy a message broker like `RabbitMQ` (or Kafka) within our Kubernetes cluster.
  This can be done using Helm charts or by defining Kubernetes Deployments and
  Services for the broker.

- **Communication via Kubernetes Services**: Our microservices will interact with
  the message broker using its Kubernetes Service name and the appropriate client
  libraries for the chosen broker (e.g., `github.com/rabbitmq/amqp091-go` for RabbitMQ).

- **Decoupled Interactions**: Services will publish messages to exchanges or
  topics hosted by the broker, and other services can subscribe to queues or
  consumer groups to process these messages asynchronously. Kubernetes ensures
  the broker is accessible within the cluster.

> !NOTE: TO know more about `topics`, `message queues` and `service brokers` read
> [Streaming Data](https://prateeksingh.tech) by me.
