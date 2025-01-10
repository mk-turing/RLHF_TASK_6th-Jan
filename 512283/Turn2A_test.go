import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
const homeUrl = 'http://localhost:8080';
const productUrl = 'http://localhost:8080/products/1';
const cartUrl = 'http://localhost:8080/cart';
const checkoutUrl = 'http://localhost:8080/checkout';

// Simulate visiting the home page
http.get(homeUrl);
sleep(1); // Simulate 1 second of browsing

// Simulate searching for a product
http.get(productUrl);
sleep(1); // Simulate 1 second of viewing product details

// Simulate adding a product to the cart
http.post(cartUrl);
sleep(1); // Simulate 1 second of cart operation

// Simulate going through checkout
http.post(checkoutUrl);
sleep(2); // Simulate 2 seconds of checkout process

check(http.get(homeUrl), {
'status was 200': (r) => r.status === 200,
});
}