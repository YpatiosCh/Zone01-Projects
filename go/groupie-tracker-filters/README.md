# 🎵 Groupie Tracker- Filters

A dynamic web application that visualizes and filters artist/band data with an elegant and responsive user interface.

## 📋 Features

### 🔍 Dynamic Search
- Real-time search functionality
- Search through artists, band members, and locations
- Instant suggestions as you type

### ⚡ Advanced Filtering System
- **Creation Date Range**: Filter bands by their formation year
- **First Album Range**: Find artists based on their first album release
- **Band Size**: Filter by number of band members
- **Locations**: Filter by concert locations with a hierarchical country/city selection
- Sliding filter panel for a clean interface
- Real-time results as you adjust filters

### 📱 Responsive Design
- Full mobile support
- Adaptive layout for different screen sizes
- Smooth transitions and animations
- Modern dark theme interface

### 🎨 UI Features
- Sliding filter panel
- Interactive range sliders
- Custom scrollbars
- Animated transitions
- Loading states
- "No results" handling
- Tooltip hints

## 🛠️ Technologies

- **Frontend**: HTML, CSS, JavaScript
- **Backend**: Go
- **Features**:
  - Custom CSS without frameworks
  - Vanilla JavaScript
  - RESTful API integration
  - Concurrent data fetching in Go

## 🏗️ Architecture

### Backend Structure
- **Models**: Data structures for artists and related information
- **Handlers**: Request processing and response generation
- **Data Service**: Concurrent API data fetching
- **Filter Package**: Data filtering and processing

### Frontend Organization
- **Components**: Modular HTML structure
- **Styling**: Organized CSS with variables
- **JavaScript**: Event handling and dynamic updates

## 💻 Installation

- Clone the repo 
```go
git clone https://platform.zone01.gr/git/ychaniot/groupie-tracker-filters 
```

## Usage 

1. Enter directory
```go
cd groupie-trackers-fitlers
```

2. Enter app directory
```go
cd src
```

3. Start the server:
```go
go run main.go
```

4. Access Application
```go
http://localhost:8080
```


## 📝 License
This project is licensed under the MIT License - see the LICENSE file for details.

## Authors
- Ypatios Chaniotakos (ychaniot)
- Marinos Kouvaras (mkouvara)




