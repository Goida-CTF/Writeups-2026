document.addEventListener('DOMContentLoaded', () => {
    console.log('Магазин "Орешник" загружен и готов к работе!');

    const productCards = document.querySelectorAll('.product-card-link');
    productCards.forEach(card => {
        card.addEventListener('mouseover', () => {
            card.style.transform = 'scale(1.05)';
            card.style.transition = 'transform 0.2s';
        });
        card.addEventListener('mouseout', () => {
            card.style.transform = 'scale(1)';
        });
    });
});
