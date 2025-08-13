import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

import { AuthProvider } from "./contexts/AuthContext";
import { DataProvider } from "./contexts/DataContext";
import { ConfirmProvider } from "./components/ConfirmDialog";
import ProtectedRoute from "./components/Auth/ProtectedRoute";

import Layout from "./components/Layout/Layout";

import Home from "./pages/Home";
import HomestayList from "./pages/HomestayList";
import About from "./pages/About";
import Login from "./pages/Login";
import Register from "./pages/Register";
import BookingHistory from "./pages/BookingHistory";
import Management from "./pages/Management";
import AddHomestay from "./pages/AddHomestay";
import EditHomestay from "./pages/EditHomestay";
import HomestayDetailManagement from "./pages/HomestayDetailManagement";
import RoomDetailPage from "./pages/RoomDetailPage";
import HomestayDetailView from "./components/Homestay/HomestayDetailView";
import GuestNewBooking from "./components/Booking/GuestNewBooking";
import VerifyAccount from "./components/Auth/VerifyAccount";

function App() {
  return (
    <ConfirmProvider>
      <Router>
        <AuthProvider>
          <DataProvider>
            <Layout>
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/homestays" element={<HomestayList />} />
                <Route path="/homestay/:id" element={<HomestayDetailView />} />
                <Route path="/about" element={<About />} />
                <Route path="/verify-account" element={<VerifyAccount />} />

                {/* Auth Routes */}
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />

                {/* Protected Routes for Guests */}
                <Route
                  path="/bookings"
                  element={
                    <ProtectedRoute requiredRoles={["guest"]}>
                      <BookingHistory />
                    </ProtectedRoute>
                  }
                />

                {/* Protected Routes for Hosts */}
                <Route
                  path="/management"
                  element={
                    <ProtectedRoute requiredRoles={["host"]}>
                      <Management />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/add-homestay"
                  element={
                    <ProtectedRoute requiredRoles={["host"]}>
                      <AddHomestay />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/management/homestay/:id"
                  element={
                    <ProtectedRoute requiredRoles={["host"]}>
                      <HomestayDetailManagement />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/management/homestay/:id/edit"
                  element={
                    <ProtectedRoute requiredRoles={["host"]}>
                      <EditHomestay />
                    </ProtectedRoute>
                  }
                />
                {/* <Route
                  path="/management/homestay/:id/rooms/add"
                  element={
                    <ProtectedRoute requiredRoles={['host']}>
                      <RoomAddPage />
                    </ProtectedRoute>
                  }
                /> */}
                <Route
                  path="/management/homestay/:homestayId/rooms/:roomId"
                  element={
                    <ProtectedRoute requiredRoles={["host"]}>
                      <RoomDetailPage />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/guest/homestay/:id/booking"
                  element={
                    <ProtectedRoute requiredRoles={["guest"]}>
                      <GuestNewBooking />
                    </ProtectedRoute>
                  }
                />
              </Routes>
            </Layout>
          </DataProvider>
        </AuthProvider>

        {/* Toast Container */}
        <ToastContainer
          position="top-right"
          autoClose={5000}
          hideProgressBar={false}
          newestOnTop={false}
          closeOnClick
          rtl={false}
          pauseOnFocusLoss
          draggable
          pauseOnHover
          theme="light"
        />
      </Router>
    </ConfirmProvider>
  );
}

export default App;
