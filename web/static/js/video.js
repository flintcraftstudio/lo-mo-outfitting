/**
 * Performance-first video loading for Lo Mo Outfitting.
 *
 * - Hero background video: two videos crossfade in a loop.
 *   Photo → river scenery → brown trout → river scenery → …
 *   Loads only on fast connections, respects reduced-motion and Save-Data.
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

  // ── Hero background video crossfade ────────────────────────────
  function initHeroVideo() {
    var video1 = document.getElementById("hero-video");
    var video2 = document.getElementById("hero-video-2");
    if (!video1) return;

    if (reducedMotion || saveData || slowConnection) {
      video1.remove();
      if (video2) video2.remove();
      return;
    }

    var videos = [video1];
    if (video2) videos.push(video2);

    // Load sources for all hero videos
    videos.forEach(function (v) {
      var source = v.querySelector("source");
      if (source && source.dataset.src) {
        source.src = source.dataset.src;
        v.load();
      }
    });

    function fadeIn(el) {
      setTimeout(function () {
        el.style.opacity = "1";
      }, 100);
    }

    function fadeOut(el) {
      el.style.opacity = "0";
    }

    if (videos.length === 1) {
      // Single video — simple loop
      video1.loop = true;
      video1.addEventListener("canplaythrough", function () {
        video1.play().then(function () { fadeIn(video1); }).catch(function () {
          video1.remove();
        });
      }, { once: true });
      return;
    }

    // Two videos — crossfade loop
    var active = 0;
    var swapping = false;

    function doSwap() {
      if (swapping) return;
      swapping = true;

      var current = videos[active];
      var next = videos[1 - active];

      next.currentTime = 0;
      next.play().then(function () {
        fadeIn(next);
        setTimeout(function () {
          fadeOut(current);
          current.pause();
          swapping = false;
        }, 1900);
        active = 1 - active;
        attachNearEnd(videos[active]);
      }).catch(function () { swapping = false; });
    }

    // Start crossfade 2s before the video ends — no gap
    function attachNearEnd(v) {
      function handler() {
        if (v === videos[active] && v.duration && v.duration - v.currentTime < 2) {
          v.removeEventListener("timeupdate", handler);
          // Pre-buffer the next video
          var next = videos[1 - active];
          next.currentTime = 0;
          next.load();
          if (next.readyState >= 3) {
            doSwap();
          } else {
            next.addEventListener("canplay", function () { doSwap(); }, { once: true });
          }
        }
      }
      v.addEventListener("timeupdate", handler);
    }

    videos.forEach(function (v) { attachNearEnd(v); });

    // Start first video once ready
    video1.addEventListener("canplaythrough", function () {
      video1.play().then(function () { fadeIn(video1); }).catch(function () {
        video1.remove();
        video2.remove();
      });
    }, { once: true });

    // Fallback for canplaythrough
    setTimeout(function () {
      if (video1.readyState < 4 && video1.parentNode) {
        video1.play().then(function () { fadeIn(video1); }).catch(function () {
          video1.remove();
          if (video2.parentNode) video2.remove();
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
