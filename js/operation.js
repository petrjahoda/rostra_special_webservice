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
    if (event.code === "Enter") {
        processOperationInput();
    }
});

function processOperationInput() {
    console.log("Operation selected: " + operationSelect.value);
    let data = {OperationSelect: operationSelect.value, OrderInput: sessionStorage.getItem("orderInput")};
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
                infoOperationInput.textContent = operationSelect.value
                sessionStorage.setItem("parovyDil", result.ParovyDil)
                infoOperationParovyDil.textContent = result.ParovyDil
                sessionStorage.setItem("seznamParovychDilu", result.SeznamParovychDilu)
                infoOperationSeznamParovychDilu.textContent = result.SeznamParovychDilu
                sessionStorage.setItem("jenPrenosMnozstvi", result.JenPrenosMnozstvi)
                infoOperationJenPrenosMnozstvi.textContent = result.JenPrenosMnozstvi
                sessionStorage.setItem("priznakMn2", result.PriznakMn2)
                infoOperationPriznakMn2.textContent = result.PriznakMn2
                sessionStorage.setItem("priznakMn3", result.PriznakMn3)
                infoOperationPriznakMn3.textContent = result.PriznakMn3
                sessionStorage.setItem("mn2Ks", result.Mn2Ks)
                infoOperationMn2Ks.textContent = result.Mn2Ks
                sessionStorage.setItem("mn3Ks", result.Mn3Ks)
                infoOperationMn3Ks.textContent = result.Mn3Ks
                sessionStorage.setItem("priznakNasobnost", result.PriznakNasobnost)
                infoOperationPriznakNasobnost.textContent = result.PriznakNasobnost
                sessionStorage.setItem("nasobnost", result.Nasobnost)
                infoOperationNasobnost.textContent = result.Nasobnost
                let pracoviste = {};
                savedWorkplaces = result.Workplaces;
                for (workplace of result.Workplaces) {
                    pracoviste[workplace.ZapsiZdroj] = workplace.ZapsiZdroj
                }
                const select = Metro.getPlugin("#workplace-select", 'select');
                select.data({
                    "Načtené pracoviště": pracoviste
                });
                workplaceSelect.focus()
            } else {
                infoError.textContent = result.OperationError;
            }
        });
    }).catch((error) => {
        infoError.textContent = error.toString()
    });
}