---
title: ethos
nav: /ethos
---

Yes, this site has an ethos.

This site was born out of the belief that the current state of the internet is wrong.
As things currently stand, the internet is predominantly controlled by big tech, but it is a decentralised system, a series of interconnected networked devices that communicate over a standardised protocol.

Users flock to singular points on the modern internet, controlled and governed by megacorps. These points are where you go to give up your privacy and digital rights, in exchange for a service.

Social media platforms offer a pretense of community. Google pretends it can help you navigate the vast web, serving you results relevant to what you're looking for, funnelling you to whoever is willing to pay the most to appear at the top of that list, and tracking you for possibilities of targeted advertising.

You're digital data is farmed and harvested for surveillance capitalism because they couldn't work out how else to squeeze out more money.

So this is my pièce de résistance. It represents my part of the internet, a small space subsumed by big tech's centralised web.

I believe that a small amount of tech literacy would go a long way. We don't need to commune on centralised platforms; rather, we can set up our own and own them.

This site runs on a single-node Kubernetes cluster hosted on an old Lenovo M910Q. It runs an `i5-6500T` (4 core) and has 8GB of RAM. 
I bought it in 2021 for £189.99.

If I check this online against AWS, an equivalent machine might be a `c6g.xlarge` (4vcpu 8GiB) with an EC2 Instance Savings Plan; for three years, it would have cost $1597.82 upfront. I've ignored the cost of electricity. I have some breathing room, but this little device has already saved me considerable money. I've used it to host quite a few services over the years such as [Tiny Tiny RSS](https://tt-rss.org/),[Miniflux](https://miniflux.app/), [Ergochat](https://ergo.chat/), [vaultwarden](https://github.com/dani-garcia/vaultwarden), [Immich](https://immich.app/), [NextCloud](https://nextcloud.com/), [Matrix](https://matrix.org/), [cigt](https://git.zx2c4.com/cgit/), etc; the list goes on. 

Currently, it only runs Miniflux, ergochat, cgit, and this site.

Someone somewhere might be angrily shaking their fist, declaring that this is no way to run a "production" environment.
I have no data backup or failover nodes — and I'd agree.
Miniflux exports its content as an XML, IRC is ephemeral, and this site is git versioned on sr.ht.

I don't need my RSS reader to be highly available; if it's broken, I'll fix it.

For things to be cost-practical, I would have to sink down to a `t2.micro` (1 core 1GiB), which costs $152 for 3 years, but then I'd only have 1 core and 1GiB.

And whilst I tend to my flock of a single node server, I learn more about the system. Kubernetes, firewalls, DNS, DHCP, etc., a responsibility delegated away or hidden behind interfaces when using cloud services.

Using your equipment is cost-effective, leaving you room to trial services you're interested in. You can also simply remove them when they're no longer required.

As for the site itself I've chosen a minimal design, keeping javascript to a minimum and serving only static content.
I want the site to work consistently across all browsers, be readable on all devices, and continue to work into the future.

The aim here is to build a place to share the things I've learnt. This space serves as an attempt to counter the low-effort writings found on Medium and write more human content as we see more AI-generated drivel. Over time, I intend to curate my posts, merging, filtering, and ultimately removing irrelevant ones.

I work as a cloud developer, writing primarily in go (this might explain my chosen minimal design choice), so the content I write through that lens.

I intend to provide information about how one can go about setting up a site like this and administering it yourself. 
As a go developer, I'll also include things I've learned and am learning.