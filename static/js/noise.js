/**
 * Add tv noise.
 *
 * @param {HTMLElement} elem
 * @param {number} alpha
 */
function addNoise(elem, alpha=255) {
  elem.id = !elem.id ? makeid(10) : elem.id;
  var styleId = `tvNoise-${elem.id}`;

  var canvas = document.createElement("canvas");
  var oldStyle = document.getElementById(styleId);
  if (oldStyle) {
    elem.removeChild(oldStyle);
  }

  var style = document.createElement("style");
  var ctx = canvas.getContext("2d");

  style.id = `tvNoise-${elem.id}`;
  style.innerText = `#${elem.id} {
    animation: tvShow 1s infinite steps(2, end);
  }`;

  let tileSize = 256;

  canvas.width = tileSize;
  canvas.height = tileSize;

  let cssKeyFrames = [];
  let frames = 16;
  let array = ctx.createImageData(tileSize, tileSize);

  for (let i = 0; i < frames; i++) {
    for (let j = 0; j < array.data.length; ) {
      const c = ~~(Math.random() * 2) * 255;
      array.data[j++] = c;
      array.data[j++] = c;
      array.data[j++] = c;
      array.data[j++] = alpha;
    }
    ctx.putImageData(array, 0, 0);

    let img = canvas.toDataURL("image/png", 0.1);
    cssKeyFrames.push(
      ~~((100 / (frames - 1)) * i) +
        "% { background: url(" +
        img +
        ") repeat top left/512px 512px; }"
    );
  }

  let styles = "@keyframes tvShow { " + cssKeyFrames.join("") + "}";
  style.innerText += styles;

  elem.appendChild(style);
}

/**
 * Remove tv noise.
 *
 * @param {HTMLElement} elem
 */
function removeNoise(elem) {
  elem.id = !elem.id ? makeid(10) : elem.id;
  var styleId = `tvNoise-${elem.id}`;
  var style = document.getElementById(styleId);
  if (style) {
    elem.removeChild(style);
  }
}

function makeid(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
