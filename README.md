## **BIlling Engine**


**Initial setup**  
please follow this document to install on local.  
https://docs.google.com/document/d/1g-NkxJF3ev6TE9ovG0zzw_gmwoL4NYL3CEpjBcXlVwM/edit?usp=sharing 

**Service Description**
Billing engine provide an API to make payment and get detail user loan outstanding. This service has scheduler that run in goroutine , to update billing statement

**API Spec**  
*Detail loan*  ( to get detail loan by cifnumber )
| method |PATH  |  Request | Response |
|--|--|--|--|
| GET | /user-outstanding?cif=*{{cif_number}}* |  | ``` {"id": 5, "user_cif": "CIF67890", "loan": 5000000, "status": true, "last_updated_at": "2025-03-06T17:47:11Z", "loan_outstanding": 5500000, "interest": 10, "is_delinquent": false}```   |  

  
*trigger job* ( this API is used to trigger job manually, this job will add new billing statement for the user loan )
| method |PATH  |  Request | Response |
|--|--|--|--|
| GET | /trigger-job |  | ``` Job triggered successfully ```   |  


*delinquents* ( to get list of delinquent user loan  )
| method |PATH  |  Request | Response |
|--|--|--|--|
| GET | /delinquents |  | ``` [{"id":5,"user_cif":"CIF67890","loan":5000000,"status":true,"last_updated_at":"2025-03-06T17:47:11Z","loan_outstanding":5500000,"interest":10,"is_delinquent":true},{"id":6,"user_cif":"CIF11223","loan":5000000,"status":true,"last_updated_at":"2025-03-06T17:47:11Z","loan_outstanding":5500000,"interest":10,"is_delinquent":true}] ```   |  


*payment* ( this API is used to do payment, user only able to pay the exact amount of payable in that week , if user has two week unpaid ,nominal payment for next payment is twice from bill payment for every week )
| method |PATH  |  Request | Response |
|--|--|--|--|
| POST | /payment | ```{"user_id": 4, "amount": 110000}``` | ``` {"code":200,"message":"Payment successful"} ```   |  



