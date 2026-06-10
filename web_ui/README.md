Put all chess piece images in the `./assets/` directory. Each file must be a `.png` and must follow the naming convention used by the engine:

```js
const imageFiles = {
  K: "wk.png", Q: "wq.png", R: "wr.png",
  B: "wb.png", N: "wn.png", P: "wp.png",
  k: "bk.png", q: "bq.png", r: "br.png",
  b: "bb.png", n: "bn.png", p: "bp.png"
};
```

Alternatively, you can change the asset location or filenames directly in `source.js` by modifying `ASSET_PATH` and/or the `imageFiles` mapping to match your preferred directory structure or naming scheme.
