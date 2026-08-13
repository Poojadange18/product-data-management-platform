const API_URL = window.location.origin;
const token = localStorage.getItem("producthub_token");
if (!token) window.location.replace("/login");

let products = [];
let pagination = { page: 1, total_pages: 1 };

async function api(path, options = {}) {
    const headers = { "Authorization": `Bearer ${localStorage.getItem("producthub_token")}`, ...(options.headers || {}) };
    const response = await fetch(`${API_URL}${path}`, { ...options, headers });
    const body = await response.json().catch(() => ({}));
    if (!response.ok) throw new Error(body.error || "Request failed");
    return body;
}

/* =========================
DOM ELEMENTS
========================= */

const productModal =
document.getElementById("productModal");

const editProductModal =
document.getElementById("editProductModal");

const productForm =
document.getElementById("productForm");

const editProductForm =
document.getElementById("editProductForm");

const productTable =
document.getElementById("productTable");

const dashboardProductTable =
document.getElementById("dashboardProductTable");

const searchInput =
document.getElementById("searchInput");

const categoryFilter =
document.getElementById("categoryFilter");

const stockFilter =
document.getElementById("stockFilter");

const sortFilter =
document.getElementById("sortFilter");

/* =========================
NAVIGATION
========================= */

const navItems =
document.querySelectorAll(".nav-item");

navItems.forEach(function (item) {


item.addEventListener("click", function (event) {

    event.preventDefault();

    const section =
        item.dataset.section;

    switchSection(section);

});


});

function switchSection(section) {


const sections = [
    "dashboard",
    "products",
    "analytics"
];

sections.forEach(function (name) {

    const element =
        document.getElementById(
            name + "Section"
        );

    if (element) {
        element.classList.remove("active");
    }

});


navItems.forEach(function (item) {

    item.classList.remove("active");

    if (
        item.dataset.section === section
    ) {
        item.classList.add("active");
    }

});


const selectedSection =
    document.getElementById(
        section + "Section"
    );

if (selectedSection) {
    selectedSection.classList.add("active");
}


const title =
    document.getElementById("pageTitle");

const subtitle =
    document.getElementById("pageSubtitle");


if (section === "dashboard") {

    title.textContent = "Dashboard";

    subtitle.textContent =
        "Overview of your product inventory";

} else if (section === "products") {

    title.textContent =
        "Products";

    subtitle.textContent =
        "Search, filter and manage products";

} else if (section === "analytics") {

    title.textContent =
        "Analytics";

    subtitle.textContent =
        "Understand your inventory performance";

}


}

/* =========================
LOAD PRODUCTS
========================= */

async function loadProducts() {


try {

    const result = await api(`/products?page=${pagination.page}&limit=20`);
    products = result.data;
    pagination = result;


    updateDashboard();

    updateCategoryFilter();

    renderProducts();

    renderDashboardProducts();

    renderAnalytics();
    renderPagination();


} catch (error) {

    console.error(
        "Error loading products:",
        error
    );


    showToast(
        "Unable to load products",
        "error"
    );

}

function renderPagination() {
    const container = document.getElementById("pagination");
    if (!container) return;
    container.innerHTML = `<button ${pagination.page <= 1 ? "disabled" : ""} onclick="changePage(-1)">Previous</button><span>Page ${pagination.page} of ${pagination.total_pages || 1} · ${pagination.total} products</span><button ${pagination.page >= pagination.total_pages ? "disabled" : ""} onclick="changePage(1)">Next</button>`;
}
function changePage(delta) { pagination.page += delta; loadProducts(); }


}

/* =========================
DASHBOARD
========================= */

function updateDashboard() {


const totalProducts =
    products.length;


let totalStock = 0;

let inventoryValue = 0;

let lowStock = 0;


products.forEach(function (product) {

    const stock =
        Number(product.stock);

    const price =
        Number(product.price);


    totalStock += stock;

    inventoryValue +=
        price * stock;


    if (stock > 0 && stock < 10) {
        lowStock++;
    }

});


document.getElementById(
    "totalProducts"
).textContent =
    totalProducts;


document.getElementById(
    "totalStock"
).textContent =
    totalStock.toLocaleString("en-IN");


document.getElementById(
    "inventoryValue"
).textContent =
    "₹" +
    inventoryValue.toLocaleString(
        "en-IN"
    );


document.getElementById(
    "lowStock"
).textContent =
    lowStock;


}

/* =========================
PRODUCT STATUS
========================= */

function getStatus(stock) {


stock = Number(stock);


if (stock === 0) {

    return {
        text: "Out of Stock",
        className: "out"
    };

}


if (stock < 10) {

    return {
        text: "Low Stock",
        className: "low"
    };

}


return {
    text: "In Stock",
    className: "healthy"
};


}

/* =========================
PRODUCT ROW
========================= */

function createProductRow(product) {


const status =
    getStatus(product.stock);


const row =
    document.createElement("tr");


row.innerHTML = `

    <td class="product-id">
        #${product.id}
    </td>

    <td class="product-name">
        ${escapeHTML(product.name)}
    </td>

    <td>
        ${escapeHTML(product.category)}
    </td>

    <td>
        ₹${Number(product.price)
            .toLocaleString("en-IN")}
    </td>

    <td>
        ${Number(product.stock)
            .toLocaleString("en-IN")}
    </td>

    <td>
        <span class="status ${status.className}">
            ${status.text}
        </span>
    </td>

    <td>

        <button
            class="action-btn edit-btn"
            onclick="editProduct(${product.id})">
            Edit
        </button>

        <button
            class="action-btn delete-btn"
            onclick="deleteProduct(${product.id})">
            Delete
        </button>

    </td>

`;


return row;


}

/* =========================
DASHBOARD TABLE
========================= */

function renderDashboardProducts() {


dashboardProductTable.innerHTML = "";


const latestProducts =
    products.slice().reverse().slice(0, 5);


if (latestProducts.length === 0) {

    dashboardProductTable.innerHTML = `
        <tr>
            <td colspan="7">
                <div class="empty-state">
                    <h3>No Products Yet</h3>
                    <p>Add your first product to get started.</p>
                </div>
            </td>
        </tr>
    `;

    return;

}


latestProducts.forEach(function (product) {

    dashboardProductTable.appendChild(
        createProductRow(product)
    );

});


}

/* =========================
FILTER PRODUCTS
========================= */

function getFilteredProducts() {


let result =
    products.slice();


const search =
    searchInput.value
        .trim()
        .toLowerCase();


const category =
    categoryFilter.value;


const stock =
    stockFilter.value;


if (search !== "") {

    result =
        result.filter(function (product) {

            return (
                product.name
                    .toLowerCase()
                    .includes(search) ||

                product.category
                    .toLowerCase()
                    .includes(search)
            );

        });

}


if (category !== "all") {

    result =
        result.filter(function (product) {

            return (
                product.category === category
            );

        });

}


if (stock === "in-stock") {

    result =
        result.filter(function (product) {

            return Number(product.stock) >= 10;

        });

}


if (stock === "low-stock") {

    result =
        result.filter(function (product) {

            const value =
                Number(product.stock);

            return value > 0 && value < 10;

        });

}


if (stock === "out-of-stock") {

    result =
        result.filter(function (product) {

            return Number(product.stock) === 0;

        });

}


return sortProducts(result);


}

/* =========================
SORT PRODUCTS
========================= */

function sortProducts(list) {


const sort =
    sortFilter.value;


if (sort === "name-asc") {

    list.sort(function (a, b) {

        return a.name.localeCompare(
            b.name
        );

    });

}


if (sort === "name-desc") {

    list.sort(function (a, b) {

        return b.name.localeCompare(
            a.name
        );

    });

}


if (sort === "price-low") {

    list.sort(function (a, b) {

        return (
            Number(a.price) -
            Number(b.price)
        );

    });

}


if (sort === "price-high") {

    list.sort(function (a, b) {

        return (
            Number(b.price) -
            Number(a.price)
        );

    });

}


if (sort === "stock-low") {

    list.sort(function (a, b) {

        return (
            Number(a.stock) -
            Number(b.stock)
        );

    });

}


if (sort === "stock-high") {

    list.sort(function (a, b) {

        return (
            Number(b.stock) -
            Number(a.stock)
        );

    });

}


return list;


}

/* =========================
PRODUCT TABLE
========================= */

function renderProducts() {


productTable.innerHTML = "";


const filtered =
    getFilteredProducts();


if (filtered.length === 0) {

    productTable.innerHTML = `
        <tr>
            <td colspan="7">
                <div class="empty-state">
                    <h3>No Products Found</h3>
                    <p>
                        Try changing your search or filters.
                    </p>
                </div>
            </td>
        </tr>
    `;

    return;

}


filtered.forEach(function (product) {

    productTable.appendChild(
        createProductRow(product)
    );

});


}

/* =========================
CATEGORY FILTER
========================= */

function updateCategoryFilter() {


const current =
    categoryFilter.value;


const categories =
    new Set();


products.forEach(function (product) {

    categories.add(
        product.category
    );

});


categoryFilter.innerHTML = `
    <option value="all">
        All Categories
    </option>
`;


Array.from(categories)
    .sort()
    .forEach(function (category) {

        const option =
            document.createElement("option");

        option.value = category;

        option.textContent = category;

        categoryFilter.appendChild(
            option
        );

    });


if (
    Array.from(categories)
        .includes(current)
) {

    categoryFilter.value =
        current;

}


}

/* =========================
SEARCH / FILTER EVENTS
========================= */

searchInput.addEventListener(
"input",
renderProducts
);

categoryFilter.addEventListener(
"change",
renderProducts
);

stockFilter.addEventListener(
"change",
renderProducts
);

sortFilter.addEventListener(
"change",
renderProducts
);

/* =========================
OPEN ADD MODAL
========================= */

function openAddModal() {


productForm.reset();

productModal.classList.add("show");


}

document.getElementById(
"dashboardAddProductBtn"
).addEventListener(
"click",
openAddModal
);

document.getElementById(
"productsAddProductBtn"
).addEventListener(
"click",
openAddModal
);

/* =========================
CLOSE ADD MODAL
========================= */

function closeAddModal() {


productModal.classList.remove(
    "show"
);

productForm.reset();


}

async function createProduct() {
    const product = {
        id: Number(document.getElementById("productId").value),
        name: document.getElementById("productName").value.trim(),
        category: document.getElementById("productCategory").value.trim(),
        price: Number(document.getElementById("productPrice").value),
        stock: Number(document.getElementById("productStock").value)
    };

    if (!product.id || !product.name || !product.category || product.price < 0 || product.stock < 0) {
        showToast("Please complete all product fields correctly", "error");
        return;
    }

    try {
        await api("/products", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(product)
        });
        closeAddModal();
        showToast("Product added successfully", "success");
        await loadProducts();
    } catch (error) {
        showToast(error.message || "Failed to add product", "error");
    }
}

document.getElementById(
"closeModalBtn"
).addEventListener(
"click",
closeAddModal
);

document.getElementById(
"cancelBtn"
).addEventListener(
"click",
closeAddModal
);

/* =========================
ADD PRODUCT
========================= */

productForm.addEventListener(
"submit",
async function (event) {


    event.preventDefault();


    const product = {

        name:
            document.getElementById(
                "productName"
            ).value.trim(),

        id: Number(document.getElementById("productId").value),

        category:
            document.getElementById(
                "productCategory"
            ).value.trim(),

        price:
            Number(
                document.getElementById(
                    "productPrice"
                ).value
            ),

        stock:
            Number(
                document.getElementById(
                    "productStock"
                ).value
            )

    };


    try {

        await api("/products", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(product) });


        closeAddModal();

        showToast(
            "Product added successfully",
            "success"
        );


        await loadProducts();


    } catch (error) {

        console.error(
            error
        );


        showToast(
        error.message,
            "error"
        );

    }

}


);

/* =========================
EDIT PRODUCT
========================= */

async function editProduct(id) {


try {

    const result = await api(`/products/${id}`);
    const product = result.data;


    document.getElementById(
        "editProductId"
    ).value = product.id;


    document.getElementById(
        "editProductName"
    ).value = product.name;


    document.getElementById(
        "editProductCategory"
    ).value = product.category;


    document.getElementById(
        "editProductPrice"
    ).value = product.price;


    document.getElementById(
        "editProductStock"
    ).value = product.stock;


    editProductModal.classList.add(
        "show"
    );


} catch (error) {

    console.error(
        error
    );


    showToast(
        "Unable to load product",
        "error"
    );

}


}

/* =========================
CLOSE EDIT MODAL
========================= */

function closeEditModal() {


editProductModal.classList.remove(
    "show"
);

editProductForm.reset();


}

document.getElementById(
"closeEditModalBtn"
).addEventListener(
"click",
closeEditModal
);

document.getElementById(
"cancelEditBtn"
).addEventListener(
"click",
closeEditModal
);

/* =========================
UPDATE PRODUCT
========================= */

editProductForm.addEventListener(
"submit",
async function (event) {


    event.preventDefault();


    const id =
        document.getElementById(
            "editProductId"
        ).value;


    const updatedProduct = {

        name:
            document.getElementById(
                "editProductName"
            ).value.trim(),

        category:
            document.getElementById(
                "editProductCategory"
            ).value.trim(),

        price:
            Number(
                document.getElementById(
                    "editProductPrice"
                ).value
            ),

        stock:
            Number(
                document.getElementById(
                    "editProductStock"
                ).value
            )

    };


    try {

        await api(`/products/${id}`, { method: "PUT", headers: { "Content-Type": "application/json" }, body: JSON.stringify(updatedProduct) });


        closeEditModal();


        showToast(
            "Product updated successfully",
            "success"
        );


        await loadProducts();


    } catch (error) {

        console.error(
            error
        );


        showToast(
            "Failed to update product",
            "error"
        );

    }

}


);

/* =========================
DELETE PRODUCT
========================= */

async function deleteProduct(id) {


const product =
    products.find(function (item) {

        return Number(item.id) === Number(id);

    });


if (!product) {
    return;
}


const confirmed =
    confirm(
        `Delete "${product.name}"?`
    );


if (!confirmed) {
    return;
}


try {

    await api(`/products/${id}`, { method: "DELETE" });


    showToast(
        "Product deleted successfully",
        "success"
    );


    await loadProducts();


} catch (error) {

    console.error(
        error
    );


    showToast(
        "Failed to delete product",
        "error"
    );

}


}

/* =========================
ANALYTICS
========================= */

function renderAnalytics() {


let totalStock = 0;

let inventoryValue = 0;

let healthy = 0;

let low = 0;

let out = 0;


products.forEach(function (product) {

    const stock =
        Number(product.stock);

    const price =
        Number(product.price);


    totalStock += stock;

    inventoryValue +=
        price * stock;


    if (stock === 0) {

        out++;

    } else if (stock < 10) {

        low++;

    } else {

        healthy++;

    }

});


document.getElementById(
    "analyticsProducts"
).textContent =
    products.length;


document.getElementById(
    "analyticsStock"
).textContent =
    totalStock.toLocaleString("en-IN");


document.getElementById(
    "analyticsValue"
).textContent =
    "₹" +
    inventoryValue.toLocaleString(
        "en-IN"
    );


document.getElementById(
    "analyticsLowStock"
).textContent =
    low;


document.getElementById(
    "healthyStockCount"
).textContent =
    healthy;


document.getElementById(
    "lowStockCountAnalytics"
).textContent =
    low;


document.getElementById(
    "outStockCount"
).textContent =
    out;


const total =
    products.length || 1;


document.getElementById(
    "healthyProgress"
).style.width =
    `${(healthy / total) * 100}%`;


document.getElementById(
    "lowProgress"
).style.width =
    `${(low / total) * 100}%`;


document.getElementById(
    "outProgress"
).style.width =
    `${(out / total) * 100}%`;


renderCategoryAnalytics();


}

/* =========================
CATEGORY ANALYTICS
========================= */

function renderCategoryAnalytics() {


const container =
    document.getElementById(
        "categoryAnalytics"
    );


container.innerHTML = "";


if (products.length === 0) {

    container.innerHTML = `
        <div class="empty-state">
            <h3>No Data</h3>
            <p>Add products to see analytics.</p>
        </div>
    `;

    return;

}


const categoryMap =
    new Map();


products.forEach(function (product) {

    const category =
        product.category;


    categoryMap.set(
        category,
        (categoryMap.get(category) || 0) + 1
    );

});


const entries =
    Array.from(
        categoryMap.entries()
    );


const maximum =
    Math.max(
        ...entries.map(
            function (entry) {
                return entry[1];
            }
        )
    );


entries
    .sort(function (a, b) {
        return b[1] - a[1];
    })
    .forEach(function (entry) {

        const category =
            entry[0];

        const count =
            entry[1];


        const percentage =
            (count / maximum) * 100;


        const item =
            document.createElement(
                "div"
            );


        item.innerHTML = `

            <div class="category-item">

                <div class="category-info">

                    <span class="category-name">
                        ${escapeHTML(category)}
                    </span>

                    <span class="category-count">
                        ${count} product${count !== 1 ? "s" : ""}
                    </span>

                </div>

            </div>

            <div class="category-bar">

                <div
                    class="category-fill"
                    style="width: ${percentage}%">
                </div>

            </div>

        `;


        container.appendChild(item);

    });


}

/* =========================
TOAST
========================= */

let toastTimer;

function showToast(message, type) {


const toast =
    document.getElementById("toast");


toast.textContent =
    message;


toast.className =
    `toast show ${type}`;


clearTimeout(toastTimer);


toastTimer =
    setTimeout(function () {

        toast.className =
            "toast";

    }, 3000);


}

/* =========================
ESCAPE HTML
========================= */

function escapeHTML(value) {


return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");


}

/* =========================
CLOSE MODALS ON OUTSIDE CLICK
========================= */

productModal.addEventListener(
"click",
function (event) {


    if (
        event.target ===
        productModal
    ) {

        closeAddModal();

    }

}


);

editProductModal.addEventListener(
"click",
function (event) {


    if (
        event.target ===
        editProductModal
    ) {

        closeEditModal();

    }

}

);

/* =========================
INITIALIZE
========================= */

loadProducts();

document.getElementById("logoutBtn").addEventListener("click", function () { localStorage.removeItem("producthub_token"); localStorage.removeItem("producthub_user"); window.location.replace("/login"); });
