# Yep

GET call example: https://yep.com/web?q=something  
API call default: https://api.yep.com/fs/2/search?client=web&gl=RS&no_correct=false&q=something&safeSearch=off&type=web  
API call load more: https://api.yep.com/fs/2/search?client=web&gl=RS&limit=31&no_correct=false&q=something&safeSearch=off&type=web

The `safeSearch` parameter can have the values: `off`, `moderate`, `strict`. Currently only `off` and `strict` are supported.  
The `type` parameter can have the values: `web`, `images`, `news`.

API logic:

1. API call without limit gives first 20 results (ranked 0-20)
2. API call for second page (limits=31) gives second 20 results (ranked 0-20)
3. API call for third page (limits=41) gives first 20 and third 20 results (0-20 are repeated of 1, 21-40 are new results)
4. API call for fourth page (limits=51) gives second 20 results (ranked 0-20 are repeated of 2, no new results, probably because it doesn't have any)
5. API call for fifth page (limits=61) gives same as 4 (probably because it doesn't have any)
