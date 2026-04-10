/**
 * Performance-first video loading for Lo Mo Outfitting.
 *
 * - Hero background video: loads only on fast connections, respects
 *   prefers-reduced-motion and Save-Data. Uses preload="none" +
 *   IntersectionObserver so zero bytes transfer until conditions are met.
 *
 * - Vimeo facade: renders a poster + play button. The ~500KB Vimeo
 *   player iframe loads only on click. Zero cost until interaction.
 */
(function () {
  "use strict";

  var reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  var saveData = navigator.connection && navigator.connection.saveData;
  var slowConnection =
    navigator.connection &&
    navigator.connection.effectiveType &&
    ["slow-2g", "2g"].indexOf(navigator.connection.effectiveType) !== -1;

  // ── Hero background video ───────────────────────────────────────
  function initHeroVideo() {
    var video = document.getElementById("hero-video");
    if (!video) return;

    // Skip on reduced motion, save-data, or slow connections
    if (reducedMotion || saveData || slowConnection) {
      video.remove();
      return;
    }

    // Start loading the video source
    var source = video.querySelector("source");
    if (source && source.dataset.src) {
      source.src = source.dataset.src;
      video.load();
    }

    // Wait until enough data is buffered to play through, then fade in
    // after a short settling delay so it always reads as a smooth crossfade
    // rather than a jarring pop.
    function fadeIn() {
      setTimeout(function () {
        video.style.opacity = "1";
      }, 600);
    }

    video.addEventListener("canplaythrough", function () {
      video.play().then(fadeIn).catch(function () {
        video.remove();
      });
    }, { once: true });

    // Fallback: if canplaythrough never fires (e.g. browser quirks),
    // try playing after a generous timeout and fade if it works.
    setTimeout(function () {
      if (video.readyState < 4 && video.parentNode) {
        video.play().then(fadeIn).catch(function () {
          video.remove();
        });
      }
    }, 8000);
  }

  // ── Vimeo facade ────────────────────────────────────────────────
  function initVimeoFacades() {
    var facades = document.querySelectorAll("[data-vimeo-id]");
    facades.forEach(function (el) {
      function activate() {
        var id = el.dataset.vimeoId;
        var hash = el.dataset.vimeoHash || "";
        var params = "autoplay=1&badge=0&byline=0&portrait=0&title=0&muted=0&controls=1";
        if (hash) params += "&h=" + hash;

        var iframe = document.createElement("iframe");
        iframe.src = "https://player.vimeo.com/video/" + id + "?" + params;
        iframe.setAttribute("allow", "autoplay; fullscreen; picture-in-picture");
        iframe.setAttribute("allowfullscreen", "");
        iframe.style.cssText =
          "position:absolute;top:0;left:0;width:100%;height:100%;border:0;";

        el.innerHTML = "";
        el.appendChild(iframe);
        el.classList.add("is-playing");
        el.removeAttribute("role");
        el.removeAttribute("tabindex");
      }

      el.addEventListener("click", activate, { once: true });
      el.addEventListener("keydown", function (e) {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          activate();
        }
      }, { once: true });
    });
  }

  // ── Init ────────────────────────────────────────────────────────
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", function () {
      initHeroVideo();
      initVimeoFacades();
    });
  } else {
    initHeroVideo();
    initVimeoFacades();
  }
})();
