---
title: 'Headers'
weight: 9
---

The `headers` config option allows you to specify the custom HTTP headers sent with each request.

```javascript {filename="Configuration"}
export const config = {
  headers: {
    "Authorization": "Bearer ey....",
    "User-Agent": "Mozilla/5.0 (Macintosh ...",
  },
  // ...
};
```

