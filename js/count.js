countBackButton.addEventListener("click", () => {
    workplaceRow.classList.remove("disabled");
    countNokInput.textContent = ""
    countOkInput.textContent = ""
    const selectNokTypes = Metro.getPlugin("#nok-type-select", 'select');
    selectNokTypes.data({});
    okRow.classList.add("disabled")
    nokRow.classList.add("disabled")
    countButton.disabled = true
    countBackButton.disabled = true
    workplaceOkButton.disabled = false
    workplaceBackButton.disabled = false
    infoRostra.textContent = ""
    workplaceSelect.focus()
})

countButton.addEventListener("click", () => {
    processCountInput();
})

countOkInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processCountInput();
    }
});

countNokInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processCountInput();
    }
});


function processCountInput() {
    console.log("Count OK entered: " + countOkInput.value);
    console.log("Count NOK entered: " + countNokInput.value);
    console.log("NOK type entered: " + nokTypeSelect.value);
    if (countOkInput.value === "") {
        countOkInput.value = "0"
        console.log("Count OK updated: " + countOkInput.value);
    }
    if (countNokInput.value === "") {
        countNokInput.value = "0"
        console.log("Count NOK updated: " + countNokInput.value);
    }
    let data = {
        WorkplaceCode: sessionStorage.getItem("workplaceCode"),
        UserId: sessionStorage.getItem("userId"),
        UserInput: sessionStorage.getItem("userInput"),
        OrderInput: sessionStorage.getItem("orderInput"),
        OperationSelect: sessionStorage.getItem("operationValue"),
        ParovyDil: sessionStorage.getItem("parovyDil"),
        SeznamParovychDilu: sessionStorage.getItem("seznamParovychDilu"),
        JenPrenosMnozstvi: sessionStorage.getItem("jenPrenosMnozstvi"),
        TypZdrojeZapsi: sessionStorage.getItem("typZdrojeZapsi"),
        ViceVp: sessionStorage.getItem("viceVp"),
        PriznakMn1: sessionStorage.getItem("priznakMn1"),
        PriznakMn2: sessionStorage.getItem("priznakMn2"),
        PriznakMn3: sessionStorage.getItem("priznakMn3"),
        Mn2Ks: sessionStorage.getItem("mn2Ks"),
        Mn3Ks: sessionStorage.getItem("mn3Ks"),
        OkCount: countOkInput.value,
        NokCount: countNokInput.value
    };
    fetch("/check_count_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                sessionStorage.setItem("okCount", countOkInput.value)
                sessionStorage.setItem("nokCount", countNokInput.value)
                sessionStorage.setItem("nokType", nokTypeSelect.value)
                okRow.classList.add("disabled")
                nokRow.classList.add("disabled")
                countButton.disabled = true
                countBackButton.disabled = true
                if (result.Transfer === "true") {
                    transferOrderButton.disabled = false
                }
                if (result.End === "true") {
                    endOrderButton.disabled = false
                }
                if (result.Clovek === "true") {
                    console.log("Clovek result: " + result.Clovek)
                    clovekRadio.disabled = false
                    clovekRadio.checked = true
                    sessionStorage.setItem("radio", "clovek")
                }
                if (result.Stroj === "true") {
                    console.log("Stroj result: " + result.Stroj)
                    strojRadio.disabled = false
                    strojRadio.checked = true
                    sessionStorage.setItem("radio", "stroj")
                }
                if (result.Serizeni === "true") {
                    console.log("Serizeni result: " + result.Serizeni)
                    serizeniRadio.disabled = false
                    serizeniRadio.checked = true
                    sessionStorage.setItem("radio", "serizeni")
                }
                infoRostra.textContent = ""
            } else {
                console.log("nok")
                infoRostra.textContent = result.CountError;
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
}