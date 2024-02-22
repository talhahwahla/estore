import React, { useState, useEffect } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';

const Estore = () => {
    const [products, setProducts] = useState([]);
    const [editingProduct, setEditingProduct] = useState(null);
    const [newProduct, setNewProduct] = useState({
        name: "",
        description: "",
        price: "",
        category: "",
    });
    const [cart, setCart] = useState([]);
    const [address, setAddress] = useState("");
    const [orderSuccess, setOrderSuccess] = useState(false); 
    const [addressError, setAddressError] = useState(false); 

    useEffect(() => {
        fetch("http://localhost:8080/")
            .then((response) => response.json())
            .then((data) => setProducts(data))
            .catch((error) => console.error("Error fetching products:", error));
    }, []);

    const isAdmin = window.location.pathname === "/admin";

    const handleEdit = (productId) => {
        const productToEdit = products.find(
            (product) => product.id === productId,
        );
        setEditingProduct(productToEdit);
    };

    const handleSave = () => {
        fetch("http://localhost:8080/admin/update", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(editingProduct),
        })
            .then((response) => {
                if (response.ok) {
                    console.log("Product updated successfully");
                    window.location.reload();
                } else {
                    console.error("Failed to update product");
                }
            })
            .catch((error) => console.error("Error updating product:", error));
    };

    const handleDelete = (productId) => {
        fetch(`http://localhost:8080/admin/delete?id=${productId}`, {
            method: "DELETE",
        })
            .then((response) => {
                if (response.ok) {
                    setProducts(
                        products.filter((product) => product.id !== productId),
                    );
                    console.log("Product deleted successfully");
                } else {
                    console.error("Failed to delete product");
                }
            })
            .catch((error) => console.error("Error deleting product:", error));
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        let newValue = value;
    
        if (name === "price" && parseFloat(value) < 1) {
            newValue = "1";
        }
    
        setEditingProduct({
            ...editingProduct,
            [name]: newValue,
        });
    };

    const handleNewProductChange = (e) => {
        const { name, value } = e.target;
        let newValue = value;

        if (name === "price" && parseFloat(value) < 1) {
            newValue = "1";
        }

        setNewProduct({
            ...newProduct,
            [name]: newValue,
        });
    };

    const handleCreate = () => {
        if (newProduct.price < 1) {
            console.error("Price must be at least 1");
            return;
        }

        fetch("http://localhost:8080/admin/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(newProduct),
        })
            .then((response) => {
                if (response.ok) {
                    console.log("Product created successfully");
                    window.location.reload();
                } else {
                    console.error("Failed to create product");
                }
            })
            .catch((error) => console.error("Error creating product:", error));
    };

    const handleAddToCart = (productId) => {
        const existingItemIndex = cart.findIndex((item) => item.id === productId);

        if (existingItemIndex !== -1) {
            
            const updatedCart = cart.map((item, index) =>
                index === existingItemIndex ? { ...item, quantity: item.quantity + 1 } : item
            );
            setCart(updatedCart);
        } else {
            
            const productToAdd = products.find((product) => product.id === productId);
            setCart([...cart, { ...productToAdd, quantity: 1 }]);
        }
    };

    const handleRemoveFromCart = (productId) => {
        setCart(cart.filter((item) => item.id !== productId));
    };

    const handleQuantityChange = (productId, newQuantity) => {
        if (newQuantity < 1) {
            
            newQuantity = 1;
        }
        setCart(
            cart.map((item) =>
                item.id === productId
                    ? { ...item, quantity: newQuantity }
                    : item,
            ),
        );
    };

    const handlePlaceOrder = () => {
        if (cart.length === 0) {
            
            return;
        }
        if (!address.trim()) {
            
            setAddressError(true);
            return;
        }
        
        setAddressError(false);

        const purchase = {
            products: cart.map((item) => ({
                product_id: item.id,
                quantity: item.quantity,
            })),
            customer_info: address,
        };
        console.log(JSON.stringify(purchase));

        fetch("http://localhost:8080/order", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(purchase),
        })
            .then((response) => {
                if (response.ok) {
                    console.log("Order placed successfully!");
                    setOrderSuccess(true); 
                    setCart([]);
                    setAddress("");
                } else {
                    console.error("Failed to place order");
                }
            })
            .catch((error) => console.error("Error placing order:", error));
    };

    
    const totalPrice = cart.reduce((acc, item) => acc + item.price * item.quantity, 0);

    return (
        <div className="container">
            <h2>Product List</h2>
            <div className="row">
                {products.map((product) => (
                    <div className="col-md-4" key={product.id}>
                        <div className="card mb-4">
                            <div className="card-body">
                                {editingProduct && editingProduct.id === product.id ? (
                                    <>
                                        <input
                                            type="text"
                                            className="form-control mb-2"
                                            name="name"
                                            value={editingProduct.name}
                                            onChange={handleInputChange}
                                        />
                                        <input
                                            type="number"
                                            className="form-control mb-2"
                                            name="price"
                                            value={editingProduct.price}
                                            onChange={handleInputChange}
                                        />
                                        <input
                                            type="text"
                                            className="form-control mb-2"
                                            name="description"
                                            value={editingProduct.description}
                                            onChange={handleInputChange}
                                        />
                                        <input
                                            type="text"
                                            className="form-control mb-2"
                                            name="category"
                                            value={editingProduct.category}
                                            onChange={handleInputChange}
                                        />
                                        {isAdmin && (
                                            <button className="btn btn-primary" onClick={handleSave}>Save</button>
                                        )}
                                    </>
                                ) : (
                                    <>
                                        <h3 className="card-title">{product.name}</h3>
                                        <p className="card-text">Price: {product.price}</p>
                                        <p className="card-text">Description: {product.description}</p>
                                        <p className="card-text">Category: {product.category}</p>
                                        {isAdmin && (
                                            <div>
                                                <button className="btn btn-secondary mr-2" onClick={() => handleEdit(product.id)}>Edit</button>
                                                <button className="btn btn-danger" onClick={() => handleDelete(product.id)}>Delete</button>
                                            </div>
                                        )}
                                        {!isAdmin && (
                                            <button className="btn btn-primary" onClick={() => handleAddToCart(product.id)}>Add to Cart</button>
                                        )}
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
            {isAdmin && (
                <div>
                    <h2>Create New Product</h2>
                    <input
                        type="text"
                        className="form-control mb-2"
                        name="name"
                        placeholder="Name"
                        value={newProduct.name}
                        onChange={handleNewProductChange}
                    />
                    <input
                        type="number"
                        className="form-control mb-2"
                        name="price"
                        placeholder="Price"
                        value={newProduct.price}
                        onChange={handleNewProductChange}
                    />
                    <input
                        type="text"
                        className="form-control mb-2"
                        name="description"
                        placeholder="Description"
                        value={newProduct.description}
                        onChange={handleNewProductChange}
                    />
                    <input
                        type="text"
                        className="form-control mb-2"
                        name="category"
                        placeholder="Category"
                        value={newProduct.category}
                        onChange={handleNewProductChange}
                    />
                    <button className="btn btn-success" onClick={handleCreate}>Create</button>
                </div>
            )}
            {!isAdmin && (
                <div>
                    <h2>Shopping Cart</h2>
                    {cart.map((item) => (
                        <div className="mb-2" key={item.id}>
                            <p>{item.name}</p>
                            <input
                                type="number"
                                className="form-control d-inline-block w-25 mr-2"
                                value={item.quantity}
                                onChange={(e) =>
                                    handleQuantityChange(
                                        item.id,
                                        parseInt(e.target.value),
                                    )
                                }
                                min="1" 
                            />
                            <button className="btn btn-danger" onClick={() => handleRemoveFromCart(item.id)}>Remove</button>
                        </div>
                    ))}
                    <p>Total Price: ${totalPrice.toFixed(2)}</p>
                    <input
                        type="text"
                        className="form-control mb-2"
                        placeholder="Address"
                        value={address}
                        onChange={(e) => setAddress(e.target.value)}
                    />
                    {addressError && (
                        <div className="alert alert-danger mt-3" role="alert">
                            Address is required.
                        </div>
                    )}
                    <button className="btn btn-primary" onClick={handlePlaceOrder} disabled={cart.length === 0}>Place Order</button>
                    {orderSuccess && (
                        <div className="alert alert-success mt-3" role="alert">
                            Order placed successfully!
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default Estore;
