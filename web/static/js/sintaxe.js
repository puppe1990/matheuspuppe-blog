/* Sintaxe journal behavior. Loaded once from the layout; every listener hangs
   off document, so a Drive morph never needs to re-attach anything (#325). */
(function () {
  "use strict";

  function on(root, event, selector, handler) {
    root.addEventListener(event, function (e) {
      var target = e.target.closest(selector);
      if (target) handler(target, e);
    });
  }

  function bodyText(key, fallback) {
    var value = document.body.dataset[key];
    return value || fallback;
  }

  function toast(message) {
    var el = document.createElement("div");
    el.className = "amarra-toast";
    el.setAttribute("role", "status");
    el.textContent = message;
    document.body.appendChild(el);
    setTimeout(function () {
      el.remove();
    }, 3000);
  }

  /* Debounced live search: the GET form serializes q + active filters and Drive
     morphs #amarra-main; the input keeps focus because idiomorph reuses the
     node with the same id. */
  var searchTimer = null;
  on(document, "input", "#search-input", function (input) {
    var form = input.form;
    if (!form) return;
    clearTimeout(searchTimer);
    searchTimer = setTimeout(function () {
      form.requestSubmit();
    }, 350);
  });

  /* Focus mode: body.amarra-focus (CSS) hides the journal chrome; Esc leaves. */
  function toggleFocus() {
    var active = document.body.classList.toggle("amarra-focus");
    document.querySelectorAll("[data-focus-toggle]").forEach(function (button) {
      button.setAttribute("aria-pressed", active ? "true" : "false");
    });
  }

  on(document, "click", "[data-focus-toggle]", function () {
    toggleFocus();
  });

  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape" && document.body.classList.contains("amarra-focus")) {
      toggleFocus();
    }
  });

  /* Font size for the reader, kept across visits. */
  var fontSize = null;
  try {
    fontSize = localStorage.getItem("sintaxe_fontsize");
  } catch (e) {}

  function applyFontSize(size) {
    var prose = document.querySelector("[data-article-prose]");
    if (!prose) return;
    prose.setAttribute("data-fontsize", size);
    document.querySelectorAll("[data-fontsize]").forEach(function (button) {
      button.setAttribute("aria-pressed", button.getAttribute("data-fontsize") === size ? "true" : "false");
    });
  }

  on(document, "click", "[data-fontsize]", function (button) {
    var size = button.getAttribute("data-fontsize");
    try {
      localStorage.setItem("sintaxe_fontsize", size);
    } catch (e) {}
    applyFontSize(size);
  });

  /* Reading progress + remaining minutes; runs once, reads the DOM each tick. */
  function updateProgress() {
    var track = document.getElementById("reading-progress");
    if (!track) return;
    var root = document.querySelector("[data-article-root]");
    if (!root) return;
    var total = document.documentElement.scrollHeight - window.innerHeight;
    var progress = 0;
    if (total > 0) {
      progress = Math.min(100, Math.max(0, Number(((window.scrollY / total) * 100).toFixed(1))));
    }
    track.style.transform = "scaleX(" + progress / 100 + ")";
    var trackHost = document.getElementById("reading-progress-track");
    if (trackHost) trackHost.setAttribute("aria-valuenow", String(Math.round(progress)));
    var readout = document.getElementById("focus-readout");
    if (readout) {
      var minutes = parseInt(root.getAttribute("data-read-minutes") || "0", 10);
      var remaining = Math.max(0, Math.ceil(minutes * (1 - progress / 100)));
      var pattern = bodyText("focusReadout", "%d%% lido · %dm restantes");
      readout.textContent = pattern.replace("%d", Math.round(progress)).replace("%d", remaining);
    }
  }

  if (document.querySelector("[data-article-root]")) {
    updateProgress();
  }
  window.addEventListener("scroll", updateProgress, { passive: true });
  window.addEventListener("resize", updateProgress, { passive: true });

  /* Share: native sheet when available, clipboard + toast otherwise. */
  on(document, "click", "[data-share]", function (button) {
    var payload = {
      title: button.getAttribute("data-share-title") || document.title,
      text: button.getAttribute("data-share-text") || "",
      url: button.getAttribute("data-share-url") || window.location.href,
    };
    if (navigator.share) {
      navigator.share(payload).catch(function () {});
      return;
    }
    if (navigator.clipboard) {
      navigator.clipboard.writeText(payload.url).then(function () {
        toast(bodyText("shareToast", "Link copiado."));
      }).catch(function () {});
    }
  });

  /* Copy buttons on code blocks. */
  on(document, "click", "[data-copy-target]", function (button) {
    var target = document.getElementById(button.getAttribute("data-copy-target"));
    if (!target || !navigator.clipboard) return;
    navigator.clipboard.writeText(target.textContent).then(function () {
      var label = button.querySelector("[data-copy-label]");
      if (!label) return;
      var original = label.textContent;
      label.textContent = bodyText("copied", "Copiado");
      setTimeout(function () {
        label.textContent = original;
      }, 2000);
    }).catch(function () {});
  });

  /* Simulated narration widget (the shipped prototype's behavior). */
  var audioTimer = null;
  var audioProgress = 0;

  function audioElements() {
    return {
      toggle: document.querySelector("[data-audio-toggle]"),
      progress: document.querySelector("[data-audio-progress]"),
      label: document.querySelector("[data-audio-label]"),
    };
  }

  on(document, "click", "[data-audio-toggle]", function () {
    var els = audioElements();
    var playing = audioTimer !== null;
    if (playing) {
      clearInterval(audioTimer);
      audioTimer = null;
      if (els.toggle) els.toggle.setAttribute("aria-pressed", "false");
      if (els.label) els.label.textContent = bodyText("audioIdle", "Ouvir");
      return;
    }
    audioTimer = setInterval(function () {
      audioProgress += 1;
      if (audioProgress > 100) {
        audioProgress = 0;
        clearInterval(audioTimer);
        audioTimer = null;
        var current = audioElements();
        if (current.toggle) current.toggle.setAttribute("aria-pressed", "false");
        if (current.label) current.label.textContent = bodyText("audioIdle", "Ouvir");
      }
      var els2 = audioElements();
      if (els2.progress) els2.progress.style.transform = "scaleX(" + audioProgress / 100 + ")";
    }, 500);
    if (els.toggle) els.toggle.setAttribute("aria-pressed", "true");
    if (els.label) els.label.textContent = bodyText("audioPlaying", "Tocando");
  });

  /* Newsletter dialog: remember where the reader was so the flash lands there. */
  on(document, "click", "[data-amarra-dialog-open]", function () {
    var next = document.querySelector("[data-newsletter-next]");
    if (next) next.value = window.location.pathname + window.location.search;
  });

  /* Metric bars ship widthless; JS paints the percentage once the DOM exists
     (and again after every Drive morph). */
  function paintBars() {
    document.querySelectorAll("[data-bar-width]").forEach(function (bar) {
      var width = parseFloat(bar.getAttribute("data-bar-width") || "0");
      bar.style.transform = "scaleX(" + Math.min(100, Math.max(0, width)) / 100 + ")";
    });
  }

  paintBars();
  document.addEventListener("amarra:morphed", function () {
    paintBars();
    applyFontSize(fontSize || "normal");
    updateProgress();
  });
  applyFontSize(fontSize || "normal");
})();
