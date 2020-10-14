const resetButton = document.getElementById("reset-button");

const userInput = document.getElementById("user-input");
const orderInput = document.getElementById("order-input");
const operationInput = document.getElementById("operation-input");
const workplaceSelect = document.getElementById("workplace-select");
const countOkInput = document.getElementById("count-ok-input");
const countNokInput = document.getElementById("count-nok-input");

const nokTypeSelect = document.getElementById("nok-type-select");

const userOkButton = document.getElementById("user-ok-button");
const orderOkButton = document.getElementById("order-ok-button");
const operationOkButton = document.getElementById("operation-ok-button");
const workplaceOkButton = document.getElementById("workplace-ok-button");
const countButton = document.getElementById("count-button");

const orderBackButton = document.getElementById("order-back-button");
const operationBackButton = document.getElementById("operation-back-button");
const workplaceBackButton = document.getElementById("workplace-back-button");
const countBackButton = document.getElementById("count-back-button");

const clovekRadio = document.getElementById("clovek-radio");
const serizeniRadio = document.getElementById("serizeni-radio");
const strojRadio = document.getElementById("stroj-radio");

const startOrderButton = document.getElementById("start-order-button");
const endOrderButton = document.getElementById("end-order-button");
const transferOrderButton = document.getElementById("transfer-order-button");

const table = document.getElementById("table");

const errorInfoPanel = document.getElementById("error-info-panel");

resetButton.addEventListener("click", () => {
    window.location.replace('');
})


// USER INPUT
userInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processUserInput();
    }
});

userOkButton.addEventListener("click", () => {
    processUserInput();
})

function processUserInput() {
    console.log("User entered: " + userInput.value);
    let data = {UserInput: userInput.value};
    fetch("/check_user_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                userInput.value = result.UserName;
                userInput.dataset.userId = result.UserId;
                userInput.dataset.userInput = result.UserInput;
                userInput.style.backgroundColor = "white"
                userInput.disabled = true;
                userOkButton.disabled = true;
                orderBackButton.disabled = false;
                orderInput.disabled = false;
                orderOkButton.disabled = false;
                orderInput.placeholder = "Zadejte číslo výrobního příkazu";
                sessionStorage.setItem("userId", result.UserId)
                sessionStorage.setItem("userName", result.UserName)
                sessionStorage.setItem("userInput", result.UserInput)
                orderInput.focus();
            } else {
                console.log(result.UserError)
                userInput.placeholder = result.UserError
                userInput.value = ""
            }
        });
    }).catch((error) => {
        errorInfoPanel.textContent = error.toString()
    });
}


// ORDER INPUT
orderBackButton.addEventListener("click", () => {
    userInput.disabled = false;
    userOkButton.disabled = false;
    orderBackButton.disabled = true;
    orderInput.disabled = true;
    orderOkButton.disabled = true;
    orderInput.placeholder = ""
    orderInput.value = ""
    userInput.value = ""
    userInput.placeholder = "Zadejte osobní číslo"
    userInput.focus()
})

orderOkButton.addEventListener("click", () => {
    processOrderInput();

})

orderInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processOrderInput();
    }
});

function processOrderInput() {
    console.log("User entered: " + orderInput.value);
    let data = {OrderInput: orderInput.value};
    fetch("/check_order_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                orderInput.value = result.OrderName;
                orderInput.dataset.userId = result.OrderId;
                orderInput.dataset.userInput = result.OrderInput;
                orderInput.style.backgroundColor = "white"
                orderInput.disabled = true;
                orderOkButton.disabled = true;
                orderBackButton.disabled = true;
                operationBackButton.disabled = false;
                operationInput.disabled = false;
                operationOkButton.disabled = false;
                operationInput.placeholder = "Zadejte číslo operace výrobního příkazu";
                sessionStorage.setItem("orderId", result.OrderId)
                sessionStorage.setItem("orderName", result.OrderName)
                sessionStorage.setItem("orderInput", result.OrderInput)
                operationInput.focus();
            } else {
                console.log(result.OrderError)
                orderInput.placeholder = result.OrderError
                orderInput.value = ""
            }
        });
    }).catch((error) => {
        errorInfoPanel.textContent = error.toString()
    });
}


// OPERATION INPUT
operationBackButton.addEventListener("click", () => {
    orderInput.disabled = false;
    orderOkButton.disabled = false;
    orderBackButton.disabled = false;
    operationBackButton.disabled = true;
    operationInput.disabled = true;
    operationOkButton.disabled = true;
    operationInput.placeholder = ""
    operationInput.value = ""
    orderInput.value = ""
    orderInput.placeholder = "Zadejte číslo výrobního příkazu";
    orderInput.focus()
})

operationOkButton.addEventListener("click", () => {
    workplaceSelect.focus()
})


// OTHER
workplaceOkButton.addEventListener("click", () => {
    countOkInput.focus()
})

countButton.addEventListener("click", () => {
    console.log("count button clicked")
})



workplaceBackButton.addEventListener("click", () => {
    operationInput.focus()
})
countBackButton.addEventListener("click", () => {
    workplaceSelect.focus()
})

