// Temporary in-memory cache for citizen photos received via DAN POST callback
const danImageCache = new Map<string, { data: string; expires: number }>();

if (typeof setInterval !== "undefined") {
  setInterval(() => {
    const now = Date.now();
    danImageCache.forEach((val, key) => {
      if (now > val.expires) danImageCache.delete(key);
    });
  }, 60_000);
}

export { danImageCache };
