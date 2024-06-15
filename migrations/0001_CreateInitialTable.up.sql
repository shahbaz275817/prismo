-- Create the Accounts table
CREATE TABLE Accounts (
                          Account_ID INT PRIMARY KEY generated always as identity,
                          Document_Number VARCHAR(15) NOT NULL
);


-- Create the OperationsTypes table
CREATE TABLE OperationsTypes (
                                 OperationType_ID INT PRIMARY KEY generated always as identity,
                                 Description VARCHAR(50) NOT NULL
);

-- Insert data into the OperationsTypes table
INSERT INTO OperationsTypes (Description) VALUES
                                                                ('Normal Purchase'),
                                                                ('Purchase with installments'),
                                                                ('Withdrawal'),
                                                                ('Credit Voucher');

-- Create the Transactions table
CREATE TABLE Transactions (
                              Transaction_ID INT PRIMARY KEY generated always as identity,
                              Account_ID INT,
                              OperationType_ID INT,
                              Amount DECIMAL(10, 2),
                              EventDate TIMESTAMP,
                              FOREIGN KEY (Account_ID) REFERENCES Accounts(Account_ID),
                              FOREIGN KEY (OperationType_ID) REFERENCES OperationsTypes(OperationType_ID)
);

