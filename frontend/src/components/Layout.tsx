import Header from './Header';
interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div className=' bg-main-black h-screen'>
      <Header />
      <main className=' px-60 pt-32'>{children}</main>
    </div>
  );
};

export default Layout;
