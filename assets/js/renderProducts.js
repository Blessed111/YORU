const productsContainer = document.querySelector('#products-container');

// Запускаем getProducts
getProducts();

// Асинхронная функция получения данных из файла products.json
async function getProducts() {
	// Получаем данные из products.json
    const response = await fetch('../assets/js/products.json');
    // Парсим данные из JSON формата в JS
    const productsArray = await response.json();
    // Запускаем ф-ю рендера (отображения товаров)
	renderProducts(productsArray);
}

function renderProducts(productsArray) {
    productsArray.forEach(function (item) {
        const productHTML = `  <div class="col-md-6">
             <div class="card mb-4" data-id="${item.id}">
              <div class="gears-card">

               <div class="card-banner">

                    <a href="#">
                      <img class="product-img" src="./assets/images/${item.imgSrc}" alt="Headphone">
                    </a>

                    <button class="share">
                      <ion-icon name="share-social"></ion-icon>
                    </button>

                    <div class="card-time-wrapper">
                      

                      <span>Best Cost</span>
                    </div>

                  </div>


              <div class="card-body text-center">

                 <div class="card-content">
                <div class="card-title-wrapper">
    
                <h4 class="item-title">${item.title}</h4>
                <p><small data-items-in-box class="text-muted">${item.itemsInBox}e-sports</small></p>
                </div>

                <div class="details-wrapper">

                  <div class="price">
                    <div class="price__weight"></div>
                    <div class="price__currency">${item.price}</div>
                  </div>
                </div>

                 <div class="details-wrapper">
                  <div class="items counter-wrapper">
                    <div class="items__control" data-action="minus">-</div>
                    <div style="color: black;" class="items__current" data-counter>1</div>
                    <div class="items__control" data-action="plus">+</div>
                  </div>
                </div>

              </div>




                  <div class="card-actions">

                    <button data-cart type="button" class="btn btn-block btn-outline-warning">Add to cart</button>

                 

                  </div>

              </div>
            </div>`;
        productsContainer.insertAdjacentHTML('beforeend', productHTML);
    });
}