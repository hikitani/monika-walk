document.addEventListener('DOMContentLoaded', () => {
    const isMobile = window.innerHeight > window.innerWidth;

    if (isMobile) {
      document.body.style.position = 'fixed';
    } else {
      document.body.style.overflow = 'hidden';
    }
})