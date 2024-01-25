function initApp() {
  const app = {
   // url: 'http://localhost:5001',
    time: null,
    activeMenu: 'pos',
    moneys: [2000, 5000, 10000, 20000, 50000, 100000],
    itemTypes: [],
    keyword: "",
    cart: [],
    orders: [],
    lineItems: [],
    cash: 0,
    change: 0,
    isProductPage: false,
    isAuthenticated: false,
    isOrderPage: false,
    isShowModalReceipt: false,
    receiptNo: null,
    receiptDate: null,
    async loadApp() {
      const response = await fetch(`http://localhost:5000/user/isCookieValid`, {
        method: 'POST',
        credentials:'include'
      })
      const data = await response.json();
      if (data.success == true) {
        this.isAuthenticated = true
        this.isProductPage = true
        this.loadProducts()
      } 
    },
   async handleLoginSubmit(e) {
      userCreds = {
        "user_name":document.getElementById("email").value,
        "password":document.getElementById("password").value
      }

      if (userCreds.user_name.length == 0 || userCreds.password.length == 0 ) {
        alert('Invalid creds!')
        return
      }
      const response = await fetch(`http://localhost:5000/user/login`, {
        method: 'POST',
        credentials: "include", //--> send/receive cookies
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(userCreds)
      })
      const data = await response.json();
      if (data.success == true) {
        this.isAuthenticated = true
        this.isProductPage = true
        this.loadProducts()
      } else {
        alert('Invalid creds!')
      }

    },
    async logout() {
      const response = await fetch(`http://localhost:5000/user/logout`, {
        method: 'POST',
        credentials:'include'
      })
      const data = await response.json();
      if (data.success == true) {
        this.isAuthenticated = false
      }
    },
    async loadProducts() {
      const response = await fetch(`http://localhost:8001/product/getproducts`)
      const data = await response.json();
      this.itemTypes = data.products;
      console.log("itemTypes loaded", this.itemTypes)
    },
    async loadOrders() {
      this.orders = [];
      const response = await fetch(`http://localhost:5001/counter/getOrders`, {
        method: 'POST',
        credentials:'include'
      })
      const res = await response.json();
      console.log("orders loaded", res.data.orders);

      this.orders = res.data.orders;
    },
    async createOrder(order) {
      const response = await fetch(`http://localhost:5001/counter/placeorder`, {
        method: 'POST',
        credentials:'include',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({"items":order})
      })
      const data = await response.json();
      console.log("orders created", data);
    },
    filteredProducts() {
      const rg = this.keyword ? new RegExp(this.keyword, "gi") : null;
      return this.itemTypes.filter((p) => !rg || p.name.match(rg));
    },
    addToCart(product) {
      const index = this.findCartIndex(product);
      if (index === -1) {
        this.cart.push({
          productType: product.type,
          image: product.image,
          name: product.name,
          price: product.price,
          qty: 1,
        });
      } else {
        this.cart[index].qty += 1;
      }
      this.beep();
      this.updateChange();
    },
    findCartIndex(product) {
      return this.cart.findIndex((p) => p.productType === product.type);
    },
    addQty(item, qty) {
      const index = this.cart.findIndex((i) => i.productType === item.productType);
      if (index === -1) {
        return;
      }
      const afterAdd = item.qty + qty;
      if (afterAdd === 0) {
        this.cart.splice(index, 1);
        this.clearSound();
      } else {
        this.cart[index].qty = afterAdd;
        this.beep();
      }
      this.updateChange();
    },
    addCash(amount) {
      this.cash = (this.cash || 0) + amount;
      this.updateChange();
      this.beep();
    },
    getItemsCount() {
      return this.cart.reduce((count, item) => count + item.qty, 0);
    },
    updateChange() {
      this.change = this.cash - this.getTotalPrice();
    },
    updateCash(value) {
      // this.cash = parseFloat(value.replace(/[^0-9]+/g, ""));
      this.cash = value;
      this.updateChange();
    },
    getTotalPrice() {
      return this.cart.reduce(
        (total, item) => total + item.qty * item.price,
        0
      );
    },
    submitable() {
      return this.change >= 0 && this.cart.length > 0;
    },
    submit() {
      const time = new Date();
      this.isShowModalReceipt = true;
      this.receiptNo = `TWPOS-KS-${Math.round(time.getTime() / 1000)}`;
      this.receiptDate = this.dateFormat(time);
    },
    closeModalReceipt() {
      this.isShowModalReceipt = false;
    },
    dateFormat(date) {
      const formatter = new Intl.DateTimeFormat('id', { dateStyle: 'short', timeStyle: 'short' });
      return formatter.format(date);
    },
    numberFormat(number) {
      // return (number || "")
      //   .toString()
      //   .replace(/^0|\./g, "")
      //   .replace(/(\d)(?=(\d{3})+(?!\d))/g, "$1.");
      return number;
    },
    priceFormat(number) {
      return number ? `${this.numberFormat(number)}$` : `0$`;
    },
    resolveImage(image) {
      return `/${image}`;
    },
    changeToProductPage() {
      this.loadProducts();
      this.isProductPage = true;
      this.isOrderPage = false;

    },
    changeToOrderPage() {
      this.loadOrders();
      this.isOrderPage = true;
      this.isProductPage = false;

    },
    clear() {
      this.cash = 0;
      this.cart = [];
      this.receiptNo = null;
      this.receiptDate = null;
      this.updateChange();
      this.clearSound();
    },
    beep() {
      this.playSound("static/sound/beep-29.mp3");
    },
    clearSound() {
      this.playSound("static/sound/button-21.mp3");
    },
    playSound(src) {
      const sound = new Audio();
      sound.src = src;
      sound.play();
      sound.onended = () => delete (sound);
    },
    printAndProceed() {
      const receiptContent = document.getElementById('receipt-content');
      const titleBefore = document.title;
      const printArea = document.getElementById('print-area');

      printArea.innerHTML = receiptContent.innerHTML;
      document.title = this.receiptNo;

      // window.print();
      this.isShowModalReceipt = false;

      printArea.innerHTML = '';
      document.title = titleBefore;

      // TODO save sale data to database

      let items = []
      for (let c of this.cart) {
        items.push({"item_type":c.productType,"quantity":c.qty})
      }

      this.createOrder(items);

      this.clear();
    }
  };

  return app;
}