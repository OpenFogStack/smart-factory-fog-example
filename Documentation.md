# Documentation

Every component logs incoming/outgoing messages and ids.
How all components are interconnected can be seen in the below figure:
![](./application.png)
It is only possible to determine the latency between components that communicate on a path without a time-based value aggregation (e.g., "emit the average every 100ms").
Four such paths exist.

## Production Path (A1)

- camera: send,camera,id                            contains a single image, emitted at fixed interval (100ms)
- check-for-defects: recv,image,id
- check-for-defects: send,cfd,id                    indicates, that a single item was defect, emitted when recv,image,id has more black than white pixels
- production-machine: recv,discard,id

## Packaging Path (A2)

- temperature-sensor: send,adapt,id                 contains a single temperature value, emitted at fixed interval (100ms)
- adapt-packaging: recv,sensor,id
OR
- production-controller: send,prodctrl,id           contains the current production rate, emitted at fixed interval (100ms)
- adapt-packaging:  recv,prodcntrl,id
THE BELOW ID CAN BE A TEMPERATURE-SENSOR ID OR AN PRODUCTION-CONTROLLER ID       
- adapt-packaging: send,adapt,id                    contains the current packaging rate and backlog, emitted on reception of recv,sensor,id or recv,prodcntrl,id
- packaging-control: recv,rate,id

## Prediction Path (A3)

- packaging-control: send,packctrl,id               contains the current packaging rate and backlog, emitted at fixed interval (100ms)
- predict-pickup: recv,input,id
- predict-pickup: send,predict,id                   contains a prediction string, emitted on reception of packaging-control data
- logistics-prognosis: recv,input,id        

## Dashboard Path (A4)

- packaging-control: send,packcntrl,id              contains the current packaging rate and backlog, emitted at fixed interval (100ms)
- aggregate: recv,data,id
- aggregate: send,aggregate,id                      contains the average packaging rate and backlog of the last x values, uses the id of the last recv,data,id message, emitted every x values (x = 10)
- generate-dashboard: recv,input,id
- generate-dashboard: send,generate_dashboard,id    contains a histogram for the average packaging rate and backlog, emitted on reception of recv,input,id
- central-office-dashboard: recv,input,id
