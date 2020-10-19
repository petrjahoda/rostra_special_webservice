countBackButton.addEventListener("click", () => {
    workplaceRow.classList.remove("disabled");
    okRow.classList.add("disabled")
    nokRow.classList.add("disabled")
    countNokInput.textContent = ""
    countOkInput.textContent = ""
    const select = Metro.getPlugin("#workplace-select", 'select');
    select.data({});
    operationSelect.focus()
})

workplaceOkButton.addEventListener("click", () => {
    processCountInput();
})


workplaceSelect.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processCountInput();
    }
});

function processCountInput() {
    console.log("Processing count input")
}