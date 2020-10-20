orderBackButton.addEventListener("click", () => {
    orderInput.placeholder = ""
    orderInput.value = ""
    userInput.value = ""
    userInput.placeholder = "Zadejte osobní číslo"
    userRow.classList.remove("disabled");
    orderRow.classList.add("disabled")
    userOkButton.disabled = false;
    orderOkButton.disabled = true;
    orderBackButton.disabled = true;
    infoRostra.textContent = ""
    infoError.textContent = ""
    infoOrderPriznakSeriovaVyroba.textContent = ""
    infoOrderInput.textContent = ""
    infoOrderName.textContent = ""
    infoOrderId.textContent = ""
    userInput.focus()
})

orderOkButton.addEventListener("click", () => {
    console.log("order button clicked")
    processOrderInput();

})

orderInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processOrderInput();
    }
});

function processOrderInput() {
    console.log("Order entered: " + orderInput.value);
    let data = {OrderInput: orderInput.value};
    fetch("/check_order_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                orderInput.value = result.OrderName;
                orderRow.classList.add("disabled");
                operationRow.classList.remove("disabled");
                orderOkButton.disabled = true;
                orderBackButton.disabled = true;
                operationOkButton.disabled = false;
                operationBackButton.disabled = false;
                operationSelect.placeholder = "Zadejte číslo operace výrobního příkazu";
                sessionStorage.setItem("orderId", result.OrderId)
                infoOrderId.textContent = result.OrderId
                sessionStorage.setItem("orderName", result.OrderName)
                infoOrderName.textContent = result.OrderName
                sessionStorage.setItem("orderInput", result.OrderInput)
                infoOrderInput.textContent = result.OrderInput
                sessionStorage.setItem("priznakSeriovaVyroba", result.PriznakSeriovaVyroba)
                infoOrderPriznakSeriovaVyroba.textContent = result.PriznakSeriovaVyroba
                let operace = {};
                for (operation of result.Operations) {
                    operace[operation.Operace] = operation.Operace + ": " + operation.Pracoviste + " [" + operation.PracovistePopis + "]"
                }
                const select = Metro.getPlugin("#operation-select", 'select');
                select.data({
                    "Načtené operace": operace
                });
                operationSelect.focus();
            } else {
                infoError.textContent = result.OrderError
                orderInput.placeholder = result.OrderError
                orderInput.value = ""
            }
        });
    }).catch((error) => {
        infoError.textContent = error.toString()
    });
}