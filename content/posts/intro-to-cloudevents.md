---
title: "Intro to CloudEvents"
date: 2020-04-13T08:30:42-07:00
draft: true
---

![CloudEvents logo](https://github.com/cncf/artwork/raw/master/projects/cloudevents/horizontal/color/cloudevents-horizontal-color.png)

I would like to do a series of CloudEvents posts on some of my takes on the
“what”, “why” and “how” but before we can get there, let’s take a moment to
understand why CloudEvents in the first place.

In the cloud space, we talk a lot about the portability of applications that
target Kubernetes. This is because you can take the manifests you develop on
your k8s cluster locally and, with a fairly high confidence level, apply them to
another cluster in any provider and they just work.

You pick VMs because you want to pick your hardware _later_.

You pick Kubernetes because you want to pick your provider _later_.

You pick CloudEvents because you want to pick your protocol _later_.

So, just like VMs and Kubernetes, CloudEvents is giving you an escape hatch from
a protocol lock-in. This becomes huge when I can write an application that runs
(and uses) the native cloud provider’s messaging systems but I do not have to
fully bake this into my application. It also means I can test locally with
something lightweight like NATS or HTTP, but ship to production with something
heavier to operate like Kafka. (I mean this in the event consumption area, not
any queue lifecycle management.)

How does CloudEvents do this? The specification defines a minimum of what is
required to route an event through a system. Just like the markings on the box
of a product from a factory: it has all the details required to let that payload
be sorted and delivered to a recipient. CloudEvents defines a base set of
attributes to stick to the outside of the box, and provisions to add custom
attributes when required. In real terms: an occurrence of an event can be
converted into a CloudEvent by picking out the required attributes and any
additional attributes that are useful.

Once you have the set of attributes for an occurance and a payload schema, you
can map these to a specific protocol as defined by the CloudEvents Bindings (or
document your proprietary custom protocol mapping). With this Binding you can
convert from the concept of the event, and the protocol version of the event. I
like to call this the canonical representation of the event.

This is the super power of CloudEvents. There is a generalized way to turn a
list of attributes and a payload into a protocol specific representation and
back. The ”and back” part is really important because it means if something
understands two or more protocols, you can create something fancy like a
protocol bridge! Or something quite complicated in practice like a function
framework that works with several protocols but the function implementers job is
trivial; they code on the canonical version of the event, never the protocol
specific version.

## The Attributes

Now let’s look at a CloudEvent and explain the reasoning behind what the
standard attributes mean.

_Required_:

- `specversion` - the version of the spec a binding must follow.
- `id` - the identifier for this event.
- `source` - the producer of the event.
- `type` - the event type.

_Optional_:

- `datacontenttype` - the content type the data should be interpreted and
  encoded as.
- `dataschema` - a uri to the schema for the data.
- `subject` - the subject of the event in the context of the event producer.
- `time` - when the occurrence happened.

Let’s break these down and think about why these mean.

`specversion` is easy, it is basically a attribute schema pointer. The spec
holds the rules on how to process an event, so this is needed to be known first
to unpack an event from a protocol.

`id` is required to allow you to understand if you have seen this event from
this producer. The spec says `id`+`source` are supposed to be globally unique.
`subject` adds an interesting element to `id`+`source` because it allows you to
keep `source` fairly high level and `subject` can point to the instance of the
producer. An example: back to the shipping boxes full of stuff; The `id`
represents the barcode on the box, it is unique. The `source` might be the
company that makes the item in the box. And the `subject` might be the warehouse
that the box came from. If you take all three of these attributes you can find
which producer made _this_ event.

`source`, `subject` are roughly clues to find the identity of producer for this
event.

The `type` could also be a kind of category or class of event. `id`+`type`
becomes the identity of the event (taking into account `source`+`subject` as
context).

`datacontenttype` tells the consumer how to unpack the payload. It could be
text, json, or even an image. The CloudEvents spec says the total size of the
event should be small though, sub 64kb on the wire, but this is a soft limit.

## Modes

What is a mode? Modes are ways to bind the attributes and payload for a
protocol. (I also might call them encodings sometimes.)

CloudEvents needed to allow for the greatest amount of coverage for the least
amount of special casing. So the easy mode for a protocol would be to support
the “**structured-mode**”: The attributes and data have been encoded into a json
object.

**Structured-mode** is not always the best choice for some protocols, as in
protocols that already have a concept of the separation between attributes and
the data. HTTP does this: headers and body. AMQP does this: headers and payload.
So CloudEvents has the concept of a “**Binary-mode**”: attributes are bound to
more native choices for a particular protocol; For example, attributes to
headers in HTTP and AMQP.

The different modes allow for more or less complex producers to emit events on a
simple to produce mode, and then downstream the middleware could be free to
convert this to another mode or protocol and retain the event integrity, meaning
the conversions are not lossy.

### Example

Let’s take a weird example. I made a tweet several days ago:

https://twitter.com/n3wscott/status/1241811694839410688

Here is _one_ way we could break this down into CloudEvent attributes:

```
Attributes,
id: 1241811694839410688
source: twitter.com
subject: n3wscott
type: status
datacontenttype: text/plain

Data,
Did I just invent the peanut butter, jelly and habanero sandwich?? 😍😍🤩🤩🔥🔥
```

Then over HTTP we can bind this as a cloudevent.

#### HTTP Structured-Mode:

```
Raw HTTP Request:
POST / HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip
Content-Length: 219
User-Agent: Go-http-client/1.1
Content-Type: application/cloudevents+json
User-Agent: Go-http-client/1.1

{
	"data": "Did I just invent the peanut butter, jelly and habanero sandwich?? 😍😍🤩🤩🔥🔥",
	"datacontenttype": "text/plain",
	"id": "1241811694839410688",
	"source": "twitter.com",
	"specversion": "1.0",
	"type": "status"
}
```

#### HTTP Binary-Mode:

```
Raw HTTP Request:
POST / HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip
Ce-Id: 1241811694839410688
Ce-Source: twitter.com
Ce-Specversion: 1.0
Ce-Type: status
Content-Length: 91
Content-Type: plain/text
User-Agent: Go-http-client/1.1

Did I just invent the peanut butter, jelly and habanero sandwich?? 😍😍🤩🤩🔥🔥
```

These two examples show what I mean by the trade offs between structured and
binary modes. On one hand, it is easy to parse the structured version. On the
other hand, I can start producing binary http cloudevents by only adding headers
to my requests in HTTP. The tradeoffs depend on your situation.

---

For more information, visit the [CloudEvents.io](https://cloudevents.io)
website.

I hope this gave some direction to why CloudEvents is a thing and why you might
care. Next up, I will look at implementing some CloudEvents producers.
