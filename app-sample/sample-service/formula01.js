var total = 0;

let ret = {
    "total": total,
    history: []
}

for (let i of $IN.items) {
    ret.history.push(`Calculating for item ${i.sku}`)
    let itemCost = i.qntd * getPrecoBySku(i.sku)
    ret.history.push(`    item cost: ${itemCost}`)
    let itemCostFrete = i.qntd * getCustoFrete($IN.from, $IN.to, i.qntd * i.peso)
    ret.history.push(`    item frete: ${itemCostFrete}`)
    let itemTaxa = getAliquotaImposto(i.sku)
    ret.history.push(`    item tax: ${itemTaxa}`)
    let custoTaxa = itemCost * itemTaxa
    ret.history.push(`    item tax cost: ${custoTaxa}`)
    ret.history.push(`item total cost: ${itemCost + itemCostFrete + custoTaxa}`)
    ret.total = ret.total + itemCost + itemCostFrete + custoTaxa;
}

ret.history.push(`GRAND TOTAL: ${ret.total}`)

// return value


ret