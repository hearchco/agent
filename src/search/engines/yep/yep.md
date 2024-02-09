# Yep

GET call example: https://yep.com/web?q=something  
API call default: https://api.yep.com/fs/2/search?client=web&gl=RS&no_correct=false&q=something&safeSearch=off&type=web  
API call load more: https://api.yep.com/fs/2/search?client=web&gl=RS&limit=31&no_correct=false&q=something&safeSearch=off&type=web

The `safeSearch` parameter can have the values: `off`, `moderate`, `strict`. Currently only `off` and `strict` are supported.  
The `type` parameter can have the values: `web`, `images`, `news`.

Unmarshalling the incoming JSON is annoying, and can probably be be optimized.