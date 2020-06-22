const express = require('express')

const app = express()
const port = 3100

app.get('/v3/transition/message', (req, res) => {
    const event = req.query.event_name

    if (event) {
        console.log("Received event " + event)
        res.sendStatus(200)
    } else {
        console.log(`Query params ${JSON.stringify(req.query)} does not contain event.`)
        res.sendStatus(400)
    }

})

app.listen(port, () => console.log(`Example app listening at http://localhost:${port}`))
