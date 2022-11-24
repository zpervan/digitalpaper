import { Link } from 'react-router-dom';

const Header = () => {
  return (
    <header className='flex px-60 py-8 justify-center border-b border-stone-800'>
      <nav className='flex flex-1 items-center text-gray-400  uppercase font-light text-base'>
        <Link to='/' className='mr-20 hover:text-main-white'>
          Home
        </Link>
        <Link to='/about' className='mr-12 hover:text-main-white'>
          About
        </Link>
      </nav>
      <h2 className='flex-2 text-main-red text-4xl font-black uppercase'>
        Digital paper
      </h2>
      <section className='flex flex-1 justify-end items-center text-gray-400  uppercase font-light text-base'>
        <Link to='/create' className='ml-12 hover:text-main-white'>
          NEW POST
        </Link>
        <p className='ml-20 hover:text-main-white'>Register</p>
      </section>
    </header>
  );
};

export default Header;
