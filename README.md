# Jaeger-Master-Slave Template

## 1. Generally used terms in jaeger

**Agent** – A network daemon that listens for spans sent over User Datagram Protocol.

**Client** – The component that implements the OpenTracing API for distributed tracing.

**Collector** – The component that receives spans and adds them into a queue to be processed.

**Console** – A UI that enables users to visualize their distributed tracing data.

**Query** – A service that fetches traces from storage.

**Span** – The logical unit of work in Jaeger, which includes the name, starting time and duration of the operation.

**Trace** – The way Jaeger presents execution requests. A trace is composed of at least one span.

## 2. Download Jaeger Docker and Run

    docker run -d --name jaeger -p 16686:16686 -p 6831:6831/udp jaegertracing/all-in-one:latest

navigate to http://localhost:16686 to access the Jaeger UI.

# Reference Document

1. https://github.com/yurishkuro/opentracing-tutorial

2. https://www.scalyr.com/blog/jaeger-tracing-tutorial
