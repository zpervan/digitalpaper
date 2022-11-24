import Layout from '../components/Layout';

const CreatePost = () => {
  return (
    <Layout>
      <main className='flex justify-center'>
        <section className='w-1/2 flex flex-col justify-center'>
          <h2 className='text-white text-3xl font-extrabold'>Add new blog</h2>
          <form className='flex flex-col mt-4'>
            <label htmlFor='text' className='text-xl text-main-white'>
              Title
            </label>
            <input
              type='text'
              className='bg-zinc-800 h-9 rounded-md pl-4 my-3'
              placeholder='Title'
            />
            <label htmlFor='textarea' className='text-xl text-main-white'>
              Content
            </label>
            <textarea
              placeholder='Here comes the content'
              className='bg-zinc-800 rounded-md px-4 py-1 my-3 h-40'
            ></textarea>
            <button
              type='submit'
              className='bg-main-red py-3 rounded-sm text-main-white text-base font-bold uppercase'
            >
              Add post
            </button>
          </form>
        </section>
      </main>
    </Layout>
  );
};

export default CreatePost;
