import Layout from '../components/Layout';

const Signup = () => {
  return (
    <Layout>
      <main className='flex justify-center'>
        <section className='w-1/2 flex flex-col justify-center'>
          <h2 className='text-white text-3xl mb-6 font-extrabold'>Sign Up</h2>
          <form className='flex flex-col mt-4 text-main-white'>
            <section className='flex justify-between '>
              <input
                type='text'
                className='bg-zinc-800 rounded-md px-4 py-4 my-3 mr-1 h-12 w-1/2'
                placeholder='First name'
              />
              <input
                type='text'
                className='bg-zinc-800 rounded-md px-4 py-4 my-3 ml-1 h-12 w-1/2'
                placeholder='Last name'
              />
            </section>
            <input
              type='text'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
              placeholder='Username'
            />
            <input
              type='email'
              placeholder='Email'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
            />
            <input
              type='password'
              placeholder='Password'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
            />
            <button
              type='submit'
              className='bg-main-red py-3 mt-4 rounded-sm text-main-white text-base font-bold uppercase'
            >
              Sign up
            </button>
          </form>
        </section>
      </main>
    </Layout>
  );
};

export default Signup;
