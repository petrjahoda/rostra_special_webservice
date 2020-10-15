operationBackButton.addEventListener("click", () => {
    orderRow.classList.remove("disabled");
    operationRow.classList.add("disabled")
    orderInput.value = ""
    orderInput.placeholder = "Zadejte číslo výrobního příkazu";
    const select = Metro.getPlugin("#operation-select", 'select');
    select.data({
    });
    orderInput.focus()
})

operationOkButton.addEventListener("click", () => {
    processOperationInput();
})

operationSelect.addEventListener("keyup", function (event) {
    console.log("key pressed")
    if (event.code === "Enter") {
        processOperationInput();
    }
});

function processOperationInput() {
    console.log("Operation selected: " + operationSelect.value);
    let data = {OperationInput: operationSelect.value, OrderInput: sessionStorage.getItem("orderInput")};
    fetch("/check_operation_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                operationRow.classList.add("disabled")
                workplaceRow.classList.remove("disabled")
                sessionStorage.setItem("operationValue", operationSelect.value)
                let pracoviste = {};
                for (workplace of result.Workplaces) {
                    pracoviste[workplace.Zapsi_zdroj] = workplace.Zapsi_zdroj
                }
                const select = Metro.getPlugin("#workplace-select", 'select');
                select.data({
                    "Načtené pracoviště": pracoviste
                });
                workplaceSelect.focus()
            } else {
                console.log(result.OperationError);
            }
        });
    }).catch((error) => {
        errorInfoPanel.textContent = error.toString()
    });

}