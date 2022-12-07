import ReactDOM from 'react-dom/client';
import App from './App';
import Post from './pages/Post';
import './index.css';
import { QueryClientProvider, QueryClient } from 'react-query';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import CreatePost from './pages/CreatePost';
import About from './pages/About';
import Signup from './pages/Signup';

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
  },
  {
    path: '/about',
    element: <About />,
  },
  {
    path: '/post/:id',
    element: <Post />,
  },
  {
    path: '/create',
    element: <CreatePost />,
  },
  {
    path: '/signup',
    element: <Signup />,
  },
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <QueryClientProvider client={queryClient}>
    <RouterProvider router={router} />
  </QueryClientProvider>
);
