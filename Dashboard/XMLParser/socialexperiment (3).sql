-- phpMyAdmin SQL Dump
-- version 4.7.4
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Feb 21, 2018 at 10:51 PM
-- Server version: 10.1.29-MariaDB
-- PHP Version: 7.2.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `socialexperiment`
--

-- --------------------------------------------------------

--
-- Table structure for table `answersdata`
--

CREATE TABLE `answersdata` (
  `Answer` varchar(255) NOT NULL,
  `Files_Name` varchar(1024) NOT NULL,
  `Latitude` varchar(1024) NOT NULL,
  `Longitude` varchar(1024) NOT NULL,
  `Question` varchar(1024) NOT NULL,
  `Type` varchar(255) NOT NULL,
  `AssignmentId` varchar(255) NOT NULL,
  `Id` int(11) NOT NULL,
  `AssetId` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



--
-- Table structure for table `asset`
--

CREATE TABLE `asset` (
  `Id` varchar(255) NOT NULL,
  `Url` varchar(255) NOT NULL,
  `Name` varchar(255) NOT NULL,
  `StartEndDestinationId` int(11) NOT NULL,
  `AssignmentId` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Table structure for table `assignment`
--

CREATE TABLE `assignment` (
  `Id` varchar(255) NOT NULL,
  `User` varchar(200) NOT NULL,
  `Project` varchar(200) NOT NULL,
  `Task` varchar(200) NOT NULL,
  `State` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Table structure for table `combinationorder`
--

CREATE TABLE `combinationorder` (
  `Id` int(11) NOT NULL,
  `COrder` int(11) NOT NULL,
  `OptionCombinationId` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `optioncombination`
--

CREATE TABLE `optioncombination` (
  `Id` int(11) NOT NULL,
  `NextQuestion` varchar(255) NOT NULL,
  `Credits` varchar(255) NOT NULL,
  `AssignmentId` varchar(255) NOT NULL,
  `AssetId` varchar(255) NOT NULL,
  `QuestionDataId` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `questiondata`
--

CREATE TABLE `questiondata` (
  `Id` int(11) NOT NULL,
  `Question` varchar(1000) NOT NULL,
  `Latitude` varchar(255) NOT NULL,
  `Longitude` varchar(255) NOT NULL,
  `Time` varchar(255) NOT NULL,
  `Frequency` varchar(255) NOT NULL,
  `Sequence` varchar(255) NOT NULL,
  `Type` varchar(255) NOT NULL,
  `Visibility` varchar(255) NOT NULL,
  `Mandatory` varchar(255) NOT NULL,
  `AssetId` varchar(255) NOT NULL,
  `QuestionId` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Table structure for table `questionoption`
--

CREATE TABLE `questionoption` (
  `Id` int(11) NOT NULL,
  `Name` varchar(255) NOT NULL,
  `NextQuestion` varchar(500) NOT NULL,
  `Credits` varchar(255) NOT NULL,
  `QuestionDataId` int(11) NOT NULL,
  `AssignmentId` varchar(255) NOT NULL,
  `AssetId` varchar(255) NOT NULL,
  `QuestionOptionId` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



--
-- Table structure for table `questionsensor`
--

CREATE TABLE `questionsensor` (
  `QuestionDataId` int(11) NOT NULL,
  `SensorId` int(11) NOT NULL,
  `Id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Table structure for table `sensor`
--

CREATE TABLE `sensor` (
  `Id` int(11) NOT NULL,
  `Name` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `sensor`
--

INSERT INTO `sensor` (`Id`, `Name`) VALUES
(1, 'Light');
INSERT INTO `sensor` (`Id`, `Name`) VALUES
(2, 'Gyroscope');
INSERT INTO `sensor` (`Id`, `Name`) VALUES
(3, 'Proximity');
INSERT INTO `sensor` (`Id`, `Name`) VALUES
(4, 'Accelerometer');
INSERT INTO `sensor` (`Id`, `Name`) VALUES
(5, 'Location');
INSERT INTO `sensor` (`Id`, `Name`) VALUES
(6, 'Noise');

-- --------------------------------------------------------

--
-- Table structure for table `sensordata`
--

CREATE TABLE `sensordata` (
  `Acceleration` text,
  `Frequency` varchar(15) DEFAULT NULL,
  `Gyroscope` text,
  `Light` text,
  `Location` text,
  `Noise` text,
  `Question` text,
  `QuestionId` text,
  `Id` int(11) NOT NULL,
  `AssignmentId` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Table structure for table `startanddestination`
--

CREATE TABLE `startanddestination` (
  `StartLatitude` varchar(255) NOT NULL,
  `StartLongitude` varchar(255) NOT NULL,
  `DestinationLatitude` varchar(255) NOT NULL,
  `DestinationLongitude` varchar(255) NOT NULL,
  `Mode` varchar(255) NOT NULL,
  `DefaultCredit` varchar(255) NOT NULL,
  `Id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `FirstName` varchar(100) NOT NULL,
  `LastName` varchar(100) NOT NULL,
  `EmailId` varchar(500) NOT NULL,
  `Password` varchar(100) NOT NULL,
  `Id` int(11) NOT NULL,
  `StudentId` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


--
-- Indexes for dumped tables
--

--
-- Indexes for table `answersdata`
--
ALTER TABLE `answersdata`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `AssignmentId` (`AssignmentId`),
  ADD KEY `AssetId` (`AssetId`);

--
-- Indexes for table `asset`
--
ALTER TABLE `asset`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `StartEndDestinationId` (`StartEndDestinationId`),
  ADD KEY `AssignmentId` (`AssignmentId`);

--
-- Indexes for table `assignment`
--
ALTER TABLE `assignment`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `combinationorder`
--
ALTER TABLE `combinationorder`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `OptionCombinationId` (`OptionCombinationId`);

--
-- Indexes for table `optioncombination`
--
ALTER TABLE `optioncombination`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `AssetId` (`AssetId`),
  ADD KEY `AssignmentId` (`AssignmentId`),
  ADD KEY `QuestionDataId` (`QuestionDataId`);

--
-- Indexes for table `questiondata`
--
ALTER TABLE `questiondata`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `AssetId` (`AssetId`);

--
-- Indexes for table `questionoption`
--
ALTER TABLE `questionoption`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `fk_question_option` (`QuestionDataId`),
  ADD KEY `AssignmentId` (`AssignmentId`),
  ADD KEY `AssetId` (`AssetId`);

--
-- Indexes for table `questionsensor`
--
ALTER TABLE `questionsensor`
  ADD PRIMARY KEY (`Id`),
  ADD KEY `QuestionDataId` (`QuestionDataId`),
  ADD KEY `SensorId` (`SensorId`);

--
-- Indexes for table `sensor`
--
ALTER TABLE `sensor`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `sensordata`
--
ALTER TABLE `sensordata`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `startanddestination`
--
ALTER TABLE `startanddestination`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`Id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `answersdata`
--
ALTER TABLE `answersdata`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `combinationorder`
--
ALTER TABLE `combinationorder`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `optioncombination`
--
ALTER TABLE `optioncombination`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `questiondata`
--
ALTER TABLE `questiondata`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `questionoption`
--
ALTER TABLE `questionoption`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `questionsensor`
--
ALTER TABLE `questionsensor`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `sensor`
--
ALTER TABLE `sensor`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `sensordata`
--
ALTER TABLE `sensordata`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=108;

--
-- AUTO_INCREMENT for table `startanddestination`
--
ALTER TABLE `startanddestination`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `answersdata`
--
ALTER TABLE `answersdata`
  ADD CONSTRAINT `answersdata_ibfk_1` FOREIGN KEY (`AssignmentId`) REFERENCES `assignment` (`Id`),
  ADD CONSTRAINT `answersdata_ibfk_2` FOREIGN KEY (`AssetId`) REFERENCES `asset` (`Id`);

--
-- Constraints for table `asset`
--
ALTER TABLE `asset`
  ADD CONSTRAINT `asset_ibfk_1` FOREIGN KEY (`StartEndDestinationId`) REFERENCES `startanddestination` (`Id`),
  ADD CONSTRAINT `asset_ibfk_2` FOREIGN KEY (`AssignmentId`) REFERENCES `assignment` (`Id`);

--
-- Constraints for table `combinationorder`
--
ALTER TABLE `combinationorder`
  ADD CONSTRAINT `combinationorder_ibfk_1` FOREIGN KEY (`OptionCombinationId`) REFERENCES `optioncombination` (`Id`);

--
-- Constraints for table `optioncombination`
--
ALTER TABLE `optioncombination`
  ADD CONSTRAINT `optioncombination_ibfk_1` FOREIGN KEY (`AssetId`) REFERENCES `asset` (`Id`),
  ADD CONSTRAINT `optioncombination_ibfk_2` FOREIGN KEY (`AssignmentId`) REFERENCES `assignment` (`Id`);

--
-- Constraints for table `questiondata`
--
ALTER TABLE `questiondata`
  ADD CONSTRAINT `questiondata_ibfk_1` FOREIGN KEY (`AssetId`) REFERENCES `asset` (`Id`);

--
-- Constraints for table `questionoption`
--
ALTER TABLE `questionoption`
  ADD CONSTRAINT `questionoption_ibfk_1` FOREIGN KEY (`AssignmentId`) REFERENCES `assignment` (`Id`),
  ADD CONSTRAINT `questionoption_ibfk_2` FOREIGN KEY (`AssetId`) REFERENCES `asset` (`Id`);

--
-- Constraints for table `questionsensor`
--
ALTER TABLE `questionsensor`
  ADD CONSTRAINT `questionsensor_ibfk_2` FOREIGN KEY (`SensorId`) REFERENCES `sensor` (`Id`);
  
  
  
  
--
-- Updating table structure for saving TimeAtSensoring and TimeAtAnswering
--

--
-- for Table answerData
-- 
ALTER TABLE `answersdata` ADD `TimeAtAnswering` DATETIME(3) NOT NULL AFTER `AssetId`;

--
-- for sensorData
-- 
ALTER TABLE `sensordata` ADD `TimeAtSensoring` DATETIME(3) NOT NULL AFTER `AssignmentId`;

--
-- for users Table
--
ALTER TABLE `users` ADD `RegistrationTime` DATETIME(3) NOT NULL AFTER `StudentId`;
ALTER TABLE `users` ADD `LastLoginTime` DATETIME(3) NOT NULL AFTER `RegistrationTime`;



-- 
-- for questionData
--
ALTER TABLE `questiondata` ADD `Vicinity` VARCHAR(10) NOT NULL AFTER `QuestionId`;

--
-- for SensorData Table
--
ALTER TABLE `sensordata` ADD `Proximity` TEXT CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL AFTER `Noise`;
ALTER TABLE `users` ADD `UserName` VARCHAR(255) NOT NULL AFTER `Password`;

COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
