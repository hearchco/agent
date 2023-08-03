# Yandex

https://yandex.com/search/?text=no way&msid=1690979305443987-8387945366794324540-balancer-l7leveler-kubr-yp-vla-97-BAL-9947&search_source=yacom_desktop_common
https://yandex.com/search/?text=a+hard+man&p=2&rnd=4498
https://yandex.com/search/?text=no way&msid=1690979305443987-8387945366794324540-balancer-l7leveler-kubr-yp-vla-97-BAL-9947&p=1

The first page has source "yacom_desktop_common", other pages not. There also seems to be a parameter `msid` for the server which receives the request. Some requests get the `rnd` parameter. The `p` parameter is for the page number, and `text` is for the search query. Ideas for purpose of [rnd](https://stackoverflow.com/questions/7821434/what-purpose-is-of-rnd-parameter-in-http-requests).

Could get the `OnPageRank` from data-cid attribute, but this would count Video and Image sections.

Interesting response headers:
+ x-yandex-items-count: 10
+ x-yandex-req-id: 1690982345124694-8192982554673989023-balancer-l7leveler-kubr-yp-vla-62-BAL-8887
+ x-yandex-sts: 1

Yandex Cookie parameter:
Cookie: yandex_gid=105422; yp=1693570508.ygu.1#4294967295.skin.s#1706747877.szm.1:1920x1080:1140x927#1693656917.csc.1#2006342320.pcs.1#1693656921.hdrc.0; yuidss=6890579431690978508; is_gdpr=0; is_gdpr_b=CNXXWRCnxgEoAg==; _yasc=VFq3QM6X57vp49OcpyOaMGoZo1Ovm3CZZl1v2aM77uFzmEg9AK6O/VpnKOQahghCUUPLr3cCG88=; i=nIPwMBx1dJXowoo3s1ZMFUkH1uIlTp+3ntYa/EZ9xXJ/8uO1NigxP0iZmT1Foselhw6yBNs5if3SRJfGHaK4drofmWI=; yandexuid=6890579431690978508; bh=EkEiR29vZ2xlIENocm9tZSI7dj0iMTEyIiwgIkNocm9taXVtIjt2PSIxMTIiLCAiTm90PUE/QnJhbmQiO3Y9IjI0IioCPzA6ByJMaW51eCI=; my=YwA=; KIykI=1; ys=wprid.1690982319460782-11608506804469128915-balancer-l7leveler-kubr-yp-vla-62-BAL-8325; bltsr=1; spravka=dD0xNjkwOTc5NjA0O2k9MTg4LjIuMjUwLjE2NTtEPTg5OUFCQzI1NDc0ODdBM0MwNDY2OTg2QkVCNTY3NTIzNDBENTc2REYxQjg0NDFDMjlERTAxNkFDOUNFMzE5QTAxMkFFRDdBQTQ0NjdGNTY5NjZENTQwRUUyQTFGQjQwMzcwODJEOTAwQTQ3QzA4RTQ4MEJBNkQzN0UwRjFEQ0YwM0U7dT0xNjkwOTc5NjA0OTkzNTQwMTE5O2g9Y2NmY2U1MTAyYzliYTRmNGE4NTY0MTU5MzMxNGFmMWQ=

base64 decoded params:
bh=A"Google Chrome";v="112", "Chromium";v="112", "Not=A?Brand";v="24"*?0:"Linux"


`spravka` parameter base64 decoded: t=1690979604;i=188.2.250.165;D=899ABC2547487A3C0466986BEB56752340D576DF1B8441C29DE016AC9CE319A012AED7AA4467F56966D540EE2A1FB4037082D900A47C08E480BA6D37E0F1DCF03E;u=1690979604993540119;h=ccfce5102c9ba4f4a85641593314af1d

t and u seem to be the unix timestamps
i seemps to be the ip
other things unclear


Look at searxng to see how they handle it - they don't see [this](https://github.com/searx/searx/issues/3210) and [this](https://github.com/searxng/searxng/issues/961).


https://yastatic.net/nearest.js probably gets the server

Keywords to look out for in javascript: 
+ `setRequestHeader`
+ `requestHeaders`
+ `cookie`
+ `initTimestamp`
+ `new Date`
+ `captcha`

Better cookie managment?: http://go-colly.org/docs/introduction/crawling/
-> https://github.com/juju/persistent-cookiejar

Browser request headers to https://yandex.com/
```
Host: yandex.com
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
DNT: 1
Connection: keep-alive
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Sec-GPC: 1
sec-ch-ua-platform: "Linux"
sec-ch-ua: "Google Chrome";v="112", "Chromium";v="112", "Not=A?Brand";v="24"
sec-ch-ua-mobile: ?0
TE: trailers
```

Response headers:
```
content-encoding: br
p3p: policyref="/w3c/p3p.xml", CP="NON DSP ADM DEV PSD IVDo OUR IND STP PHY PRE NAV UNI"
report-to: { "group": "network-errors", "max_age": 100, "endpoints": [{"url": "https://dr.yandex.net/nel", "priority": 1}, {"url": "https://dr2.yandex.net/nel", "priority": 2}]}
reporting-endpoints: default="https://yandex.com/portal/front/reports/?slots=681853%2C0%2C16&region=105422&reqid=1690993131724565-1610792554290551541-balancer-l7leveler-kubr-yp-vla-78-BAL-3438&dc=sas&page=desktop.global&enableOtherTypes=0"
cache-control: no-cache,no-store,max-age=0,must-revalidate
x-yandex-req-id: 1690993131724565-1610792554290551541-balancer-l7leveler-kubr-yp-vla-78-BAL-3438
last-modified: Wed, 02 Aug 2023 16:18:51 GMT
nel: {"report_to": "network-errors", "max_age": 100, "success_fraction": 0.001, "failure_fraction": 0.1}
date: Wed, 02 Aug 2023 16:18:51 GMT
set-cookie: yandex_gid=105422; Path=/; Domain=yandex.com; Expires=Fri, 01 Sep 2023 16:18:51 GMT; Secure; SameSite=None
set-cookie: yp=1693585131.ygu.1#4294967295.skin.s; Path=/; Domain=yandex.com; Expires=Sat, 30 Jul 2033 16:18:51 GMT; Secure; SameSite=None
set-cookie: yuidss=7905515411690993131; Path=/; Domain=yandex.com; Expires=Sat, 30 Jul 2033 16:18:51 GMT; Secure; SameSite=None
set-cookie: is_gdpr=0; Path=/; Domain=.yandex.com; Expires=Fri, 01 Aug 2025 16:18:51 GMT; SameSite=None; Secure
set-cookie: is_gdpr_b=CNXXWRCrxgEoAg==; Path=/; Domain=.yandex.com; Expires=Fri, 01 Aug 2025 16:18:51 GMT; SameSite=None; Secure
set-cookie: _yasc=iPjny4TWibiFlemgHwa9WGWSq+fhU+MX/YCl2Trl9LrlHVlMSxLOKpHJgKVGhw==; domain=.yandex.com; path=/; expires=Sat, 30 Jul 2033 16:18:51 GMT; secure
set-cookie: i=IamW/DBcSMXyADQFc4lSE9eKnTYFEs05XKeaKMjSOeWRS3E7PcmRGlS7bQQBDTUR5sTBeKW8B9DVWB5xzucRFzm0QFA=; Expires=Fri, 01-Aug-2025 16:18:51 GMT; Domain=.yandex.com; Path=/; Secure; HttpOnly; SameSite=None
set-cookie: yandexuid=7905515411690993131; Expires=Fri, 01-Aug-2025 16:18:51 GMT; Domain=.yandex.com; Path=/; Secure; SameSite=None
set-cookie: bh=EkEiR29vZ2xlIENocm9tZSI7dj0iMTEyIiwgIkNocm9taXVtIjt2PSIxMTIiLCAiTm90PUE/QnJhbmQiO3Y9IjI0IioCPzA6ByJMaW51eCI=; Expires=Thu, 01-Aug-2024 16:18:51 GMT; Path=/
vary: Cookie,Accept-Language,Accept-Encoding
content-security-policy: report-uri https://csp.yandex.net/csp?project=morda&from=morda.big.com&showid=1690993131724565-1610792554290551541-balancer-l7leveler-kubr-yp-vla-78-BAL-3438&h=stable-portal-mordago-158.sas.yp-c.yandex.net&yandexuid=7905515411690993131&&version=2023-08-01-338.2&adb=0;connect-src *.strm.yandex.net mc.yandex.com yandex.com yastatic.net yastat.net mc.yandex.ru *.mc.yandex.ru adstat.yandex.ru mc.admetrica.ru;img-src *.verify.yandex.ru *.ya.ru *.yandex.ru ya.ru yabs.yandex.by yabs.yandex.kz yabs.yandex.ru yabs.yandex.uz yandex.ru 'self' yastatic.net data: yandex.com favicon.yandex.net avatars.mds.yandex.net mc.admetrica.ru mc.yandex.com *.mc.yandex.ru adstat.yandex.ru mc.yandex.ru;script-src 'nonce-rPCFMBgYZRmabV1akljnqQ==' mc.yandex.com yastatic.net yandex.com mc.yandex.ru *.mc.yandex.ru adstat.yandex.ru;child-src *.ya.ru *.yandex.ru ya.ru yandex.ru yastatic.net yandex.com mc.yandex.ru mc.yandex.md mc.yandex.com *.ya.ru *.yandex.ru ya.ru yandex.ru;style-src 'unsafe-inline' yastatic.net;default-src yastatic.net yastat.net;font-src yastatic.net
content-type: text/html; charset=UTF-8
x-frame-options: DENY
x-content-type-options: nosniff
strict-transport-security: max-age=31536000; includeSubDomains
accept-ch: Sec-CH-UA-Platform-Version, Sec-CH-UA-Mobile, Sec-CH-UA-Model, Sec-CH-UA, Sec-CH-UA-Full-Version-List, Sec-CH-UA-WoW64, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Platform, Sec-CH-UA-Full-Version, Viewport-Width, DPR, Device-Memory, RTT, Downlink, ECT
expires: Wed, 02 Aug 2023 16:18:51 GMT
link: <//yastatic.net/jquery/2.1.4/jquery.min.js>;  rel="preload"; as="script"; crossorigin="anonymous";
link: <https://yastatic.net/s3/home-static/_/_/3/W2txnAg-OgYrlvrPnhKuvyBvQ.js>;  rel="preload"; as="script"; crossorigin="anonymous";
X-Firefox-Spdy: h2
```

Then, search query: https://yandex.com/search/?text=some query&msid=1690993131724565-1610792554290551541-balancer-l7leveler-kubr-yp-vla-78-BAL-3438&search_source=yacom_desktop_common

Request headers:
```
Host: yandex.com
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
Referer: https://yandex.com/
DNT: 1
Connection: keep-alive
Cookie: yandex_gid=105422; yp=1693585131.ygu.1#4294967295.skin.s#1706761134.szm.1:1920x1080:879x927#1693671546.csc.1; yuidss=7905515411690993131; is_gdpr=0; is_gdpr_b=CNXXWRCrxgEoAg==; _yasc=7Qix0eMxgPdjVhsBWnFoO6BcXcyp3OiKvbu+R/Bk/cCGC72DQm0By3Ve6xW+4L1D; i=IamW/DBcSMXyADQFc4lSE9eKnTYFEs05XKeaKMjSOeWRS3E7PcmRGlS7bQQBDTUR5sTBeKW8B9DVWB5xzucRFzm0QFA=; yandexuid=7905515411690993131; bh=EkEiR29vZ2xlIENocm9tZSI7dj0iMTEyIiwgIkNocm9taXVtIjt2PSIxMTIiLCAiTm90PUE/QnJhbmQiO3Y9IjI0IioCPzA6ByJMaW51eCI=; KIykI=1; my=YwA=
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: same-origin
Sec-Fetch-User: ?1
Sec-GPC: 1
sec-ch-ua-platform: "Linux"
sec-ch-ua: "Google Chrome";v="112", "Chromium";v="112", "Not=A?Brand";v="24"
sec-ch-ua-mobile: ?0
TE: trailers
```

Response headers:
```
content-security-policy: child-src 'self' data: blob: yabrowser: yandexadexchange.net *.yandexadexchange.net *.kinopoisk.ru www.youtube.com video.khl.ru www.video.khl.ru api-video.khl.ru www.api-video.khl.ru 1tv.ru www.1tv.ru stream.1tv.ru www.stream.1tv.ru player.vgtrk.com www.player.vgtrk.com my.ntv.ru www.my.ntv.ru www.ntv.ru otr.webcaster.pro www.otr.webcaster.pro news.sportbox.ru yabs.yandex.ru paymentcard.yamoney.ru widget.bookform.ru *.turbopages.org *.turbo.site mc.yandex.md yoomoney.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;connect-src 'self' blob: wss://*.yandex.net wss://yandex.net wss://*.yandex.ru wss://yandex.ru wss://*.yandex.com wss://yandex.com yandexmetrica.com:* mc.admetrica.ru mc.yandex.md clck.ru an.yandex.ru jstracer.yandex.ru amc.yandex.ru strm.yandex.ru *.strm.yandex.ru *.strm.yandex.net verify.yandex.ru *.verify.yandex.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru mc.yandex.com;default-src 'self' blob: yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;img-src * data: blob:;script-src 'self' 'unsafe-inline' 'unsafe-eval' yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru 'nonce-6906' mc.yandex.com;style-src blob: 'self' 'unsafe-inline' yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;frame-src 'self' data: blob: yabrowser: yandexadexchange.net *.yandexadexchange.net *.kinopoisk.ru www.youtube.com video.khl.ru www.video.khl.ru api-video.khl.ru www.api-video.khl.ru 1tv.ru www.1tv.ru stream.1tv.ru www.stream.1tv.ru player.vgtrk.com www.player.vgtrk.com my.ntv.ru www.my.ntv.ru www.ntv.ru otr.webcaster.pro www.otr.webcaster.pro news.sportbox.ru yabs.yandex.ru paymentcard.yamoney.ru widget.bookform.ru *.turbopages.org *.turbo.site mc.yandex.md yoomoney.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru mc.yandex.com;media-src 'self' data: blob: *.s3.dzeninfra.ru strm.yandex.ru *.strm.yandex.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;font-src 'self' data: blob: yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;frame-ancestors yandex.com *.yandex.com yandex.ru *.yandex.ru sandbox.toloka.yandex.com sandbox.iframe-toloka.com iframe-toloka.com yang.yandex-team.ru;report-uri https://csp.yandex.net/csp?from=web4%3Adesktop&project=web4&reqid=1690993398159066-2398882240994117353-vla1-2557-vla-l7-balancer-exp-8080-BAL-5672&yandexuid=7905515411690993131&yandex_login=;
nel: {"report_to": "network-errors", "max_age": 100, "success_fraction": 0.001, "failure_fraction": 0.1}
x-content-type-options: nosniff
set-cookie: ys=wprid.1690993398159066-2398882240994117353-vla1-2557-vla-l7-balancer-exp-8080-BAL-5672; path=/; domain=.yandex.com; SameSite=None; Secure
set-cookie: _yasc=IWJGpqhdJGBVR3qrK6kqie4k8LxVZQZ9ryWB0CjNB6dqNKCMvBM1TTNjfqAy6vB1B8E=; domain=.yandex.com; path=/; expires=Sat, 30 Jul 2033 16:23:18 GMT; secure
x-frame-options: DENY
x-yandex-sts: 1
expires: Wed, 02 Aug 2023 16:28:18 GMT
expires: Wed, 02 Aug 2023 16:28:18 GMT
x-yandex-req-id: 1690993398159066-2398882240994117353-vla1-2557-vla-l7-balancer-exp-8080-BAL-5672
accept-ch: Sec-CH-UA-Platform-Version, Sec-CH-UA-Mobile, Sec-CH-UA-Model, Sec-CH-UA, Sec-CH-UA-Full-Version-List, Sec-CH-UA-WoW64, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Platform, Sec-CH-UA-Full-Version, Viewport-Width, DPR, Device-Memory, RTT, Downlink, ECT
x-yandex-items-count: 10
report-to: { "group": "network-errors", "max_age": 100, "endpoints": [{"url": "https://dr.yandex.net/nel", "priority": 1}, {"url": "https://dr2.yandex.net/nel", "priority": 2}]}
cache-control: private
cache-control: private, max-age=300
cache-control: private, max-age=300, no-transform
content-encoding: br
content-type: text/html; charset=utf-8
strict-transport-security: max-age=31536000; includeSubDomains
X-Firefox-Spdy: h2
```

Second page: https://yandex.com/search/?text=some query&msid=1690993131724565-1610792554290551541-balancer-l7leveler-kubr-yp-vla-78-BAL-3438&cee=1&p=1

Request Headers:
```
Host: yandex.com
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
Referer: https://yandex.com/
DNT: 1
Connection: keep-alive
Cookie: yandex_gid=105422; yp=1693585131.ygu.1#4294967295.skin.s#1706761134.szm.1%3A1920x1080%3A879x927#1693671546.csc.1#2006353399.pcs.1#1693671800.hdrc.0; yuidss=7905515411690993131; is_gdpr=0; is_gdpr_b=CNXXWRCrxgEoAg==; _yasc=IWJGpqhdJGBVR3qrK6kqie4k8LxVZQZ9ryWB0CjNB6dqNKCMvBM1TTNjfqAy6vB1B8E=; i=IamW/DBcSMXyADQFc4lSE9eKnTYFEs05XKeaKMjSOeWRS3E7PcmRGlS7bQQBDTUR5sTBeKW8B9DVWB5xzucRFzm0QFA=; yandexuid=7905515411690993131; bh=EkEiR29vZ2xlIENocm9tZSI7dj0iMTEyIiwgIkNocm9taXVtIjt2PSIxMTIiLCAiTm90PUE/QnJhbmQiO3Y9IjI0IioCPzA6ByJMaW51eCI=; KIykI=1; my=YwA=; ys=wprid.1690993398159066-2398882240994117353-vla1-2557-vla-l7-balancer-exp-8080-BAL-5672; bltsr=1
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: same-origin
Sec-Fetch-User: ?1
Sec-GPC: 1
sec-ch-ua-platform: "Linux"
sec-ch-ua: "Google Chrome";v="112", "Chromium";v="112", "Not=A?Brand";v="24"
sec-ch-ua-mobile: ?0
TE: trailers
```

Response Headers:
```
content-security-policy: child-src 'self' data: blob: yabrowser: yandexadexchange.net *.yandexadexchange.net *.kinopoisk.ru www.youtube.com video.khl.ru www.video.khl.ru api-video.khl.ru www.api-video.khl.ru 1tv.ru www.1tv.ru stream.1tv.ru www.stream.1tv.ru player.vgtrk.com www.player.vgtrk.com my.ntv.ru www.my.ntv.ru www.ntv.ru otr.webcaster.pro www.otr.webcaster.pro news.sportbox.ru yabs.yandex.ru paymentcard.yamoney.ru widget.bookform.ru *.turbopages.org *.turbo.site mc.yandex.md yoomoney.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;connect-src 'self' blob: wss://*.yandex.net wss://yandex.net wss://*.yandex.ru wss://yandex.ru wss://*.yandex.com wss://yandex.com yandexmetrica.com:* mc.admetrica.ru mc.yandex.md clck.ru an.yandex.ru jstracer.yandex.ru amc.yandex.ru strm.yandex.ru *.strm.yandex.ru *.strm.yandex.net verify.yandex.ru *.verify.yandex.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru mc.yandex.com;default-src 'self' blob: yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;img-src * data: blob:;script-src 'self' 'unsafe-inline' 'unsafe-eval' yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru 'nonce-10606' mc.yandex.com;style-src blob: 'self' 'unsafe-inline' yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;frame-src 'self' data: blob: yabrowser: yandexadexchange.net *.yandexadexchange.net *.kinopoisk.ru www.youtube.com video.khl.ru www.video.khl.ru api-video.khl.ru www.api-video.khl.ru 1tv.ru www.1tv.ru stream.1tv.ru www.stream.1tv.ru player.vgtrk.com www.player.vgtrk.com my.ntv.ru www.my.ntv.ru www.ntv.ru otr.webcaster.pro www.otr.webcaster.pro news.sportbox.ru yabs.yandex.ru paymentcard.yamoney.ru widget.bookform.ru *.turbopages.org *.turbo.site mc.yandex.md yoomoney.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru mc.yandex.com;media-src 'self' data: blob: *.s3.dzeninfra.ru strm.yandex.ru *.strm.yandex.ru yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;font-src 'self' data: blob: yandex.com *.yandex.com yandex.ru *.yandex.ru yastat.net yastatic.net *.yastatic.net yandex.net *.yandex.net ya.ru *.ya.ru;frame-ancestors yandex.com *.yandex.com yandex.ru *.yandex.ru sandbox.toloka.yandex.com sandbox.iframe-toloka.com iframe-toloka.com yang.yandex-team.ru;report-uri https://csp.yandex.net/csp?from=web4%3Adesktop&project=web4&reqid=1690993495067326-17948815522177371521-vla1-2557-vla-l7-balancer-exp-8080-BAL-9372&yandexuid=7905515411690993131&yandex_login=;
nel: {"report_to": "network-errors", "max_age": 100, "success_fraction": 0.001, "failure_fraction": 0.1}
x-content-type-options: nosniff
set-cookie: ys=wprid.1690993495067326-17948815522177371521-vla1-2557-vla-l7-balancer-exp-8080-BAL-9372; path=/; domain=.yandex.com; SameSite=None; Secure
x-frame-options: DENY
x-yandex-sts: 1
expires: Wed, 02 Aug 2023 16:29:55 GMT
expires: Wed, 02 Aug 2023 16:29:55 GMT
x-yandex-req-id: 1690993495067326-17948815522177371521-vla1-2557-vla-l7-balancer-exp-8080-BAL-9372
accept-ch: Sec-CH-UA-Platform-Version, Sec-CH-UA-Mobile, Sec-CH-UA-Model, Sec-CH-UA, Sec-CH-UA-Full-Version-List, Sec-CH-UA-WoW64, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Platform, Sec-CH-UA-Full-Version, Viewport-Width, DPR, Device-Memory, RTT, Downlink, ECT
x-yandex-items-count: 10
report-to: { "group": "network-errors", "max_age": 100, "endpoints": [{"url": "https://dr.yandex.net/nel", "priority": 1}, {"url": "https://dr2.yandex.net/nel", "priority": 2}]}
cache-control: private
cache-control: private, max-age=300
cache-control: private, max-age=300, no-transform
content-encoding: br
content-type: text/html; charset=utf-8
strict-transport-security: max-age=31536000; includeSubDomains
X-Firefox-Spdy: h2
```

**Response is not actually time based.**

The only important header is `Cookie:` - the other ones don't matter. It seems that almost every time you send a request without cookies, you will be hit with a capcha. Yandex makes [their own captcha](https://cloud.yandex.com/en/services/smartcaptcha), maybe try to crack it?

From [support](https://yandex.com/support/smart-captcha/): 
> A Yandex service may be blocked if it receives many similar requests from users or programs. For example, this may happen if multiple people use a Yandex service from devices that are connected to the internet from the same IP address. In this case, Yandex interprets them all as one user and asks for additional confirmation.

Figuring out the `spravka` parameter is probably the key to success.
Sending identical requests is /VERY/ bad/obvious.

![repeated](site/image.png)