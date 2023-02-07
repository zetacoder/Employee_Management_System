# Employee_Management_System
Simple app that combinates an API REST architecture with CRUD principles using an MySQL database.
Allows to manage your employees, hire nee employees, retrieve their information and if they doesnt work enough they can be fired.

## SERVER, ROUTES AND FUNCTIONALITIES
#### Server: localhost:8080

#### Routes and functionalities:

* / -> **GET** all the employees
* /employee/:id -> **GET** some employee by id. E.g: localhost:8080/3
* /employee      -> **POST** creates newemployee. Requiered fields: Name, LastName, JobTitle, SalarabyDay, WorkedDays, ExpectedWorkedDays
* /mustbefired -> **GET** shows all employees that must be fired. This, if they expected worked days vs real worked days difference is > 3.
* /mustbefired -> **DELETE** erase all employees that must be fired from the database.

## EXAMPLES:

**1. CREATE NEW EMPLOYEES**
Two employees added to the database.
![image](https://user-images.githubusercontent.com/71451124/217386528-64a3a6be-a373-4ae9-946f-89b2d4202ed4.png)
![image](https://user-images.githubusercontent.com/71451124/217386787-2da7f638-9d06-4611-bb23-b997c5cf55bd.png)

**2. SHOW ALL EMPLOYEES**
![image](https://user-images.githubusercontent.com/71451124/217387038-94aca585-7757-4f80-b6c2-3c360436bfcf.png)

Retrieving the information, we can see that they were uploaded into the database:
![image](https://user-images.githubusercontent.com/71451124/217387386-897ccf73-f14c-4cc9-ae5d-d8dd517c3827.png)

**3. SHOW ALL THE EMPLOYEES IN CONDITION TO BE FIRED**
If the difference between the real worked days and expected are more than 3, that employee should be fired.
![image](https://user-images.githubusercontent.com/71451124/217387707-4e0e817e-5b1a-4110-98fb-5cd9c69687d6.png)
In this case the differece was more than 3, so it should be fired. Check previous examples to look the worked hours.

**4. FIRE ALL EMPLOYEES THAT SHOULD BE FIRED**
If met the criteria, would be deleted from the database (and company).
![image](https://user-images.githubusercontent.com/71451124/217388153-6409a99c-14cb-4c11-8dff-5fee8839e25f.png)
