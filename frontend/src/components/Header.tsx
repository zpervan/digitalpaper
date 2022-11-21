const Header = () => {
  return (
    <header className='flex px-60 py-8 justify-center'>
      <nav className='flex flex-1 items-center text-gray-400  uppercase font-light text-sm'>
        <a href='/' className='mr-8 hover:text-main-white'>
          Home
        </a>
        <a href='/about' className='mr-8 hover:text-main-white'>
          About
        </a>
      </nav>
      <h2 className='flex-2 text-main-red text-4xl font-black uppercase'>
        Digital paper
      </h2>
      <section className='flex flex-1 justify-end items-center text-gray-400  uppercase font-light text-sm'>
        <p className='ml-8 hover:text-main-white'>NEW POST</p>
        <p className='ml-8 hover:text-main-white'>Register</p>
      </section>
    </header>
  );
};

export default Header;
